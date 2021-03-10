package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"

	_ "net/http/pprof"

	"github.com/bartOssh/go_basilisk/debugger"
	_ "github.com/bartOssh/go_basilisk/docs"
	services "github.com/bartOssh/go_basilisk/services"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	srvReadWriteTimeout = 15
	srvIdleTimeout      = 60
	exitTimeout         = 30
)

var (
	srvAddressAndPort string
	dbClient          services.TokenAcctions
	appToken          string
	isDebug           bool
)

// URLRequestBody describes schema of request body for web page url to screenshot
type URLRequestBody struct {
	URL string `json:"url"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("cannot load .env file %s\n looking for global variables", err)
	}
	isDebug = os.Getenv("DEBUG") == "true"
	srvAddressAndPort = os.Getenv("SERVER_IP_PORT")
	dbURI := os.Getenv("MONGODB_ADDON_URI")
	dbName := os.Getenv("MONGODB_ADDON_DB")
	setToken := os.Getenv("SET_RESET_TOKEN")
	dbClient, err = services.NewMongoClient(dbURI, dbName)
	defer dbClient.Close()
	if setToken == "true" {
		err = dbClient.SetToken()
		if err != nil {
			log.Fatalf("cannot set app token %s", err)
		}
	}
	appToken, err = dbClient.GetToken()
	if err != nil {
		log.Fatalf("cannot read app token %s", err)
	}
	log.Printf("app token set to %s, note it down\n", appToken)
	if err != nil {
		log.Fatalf("cannot connect to data base, error: %s", err)
	}
	log.Println("initialization of API Client with success")
}

// @title Go Basilisk
// @version 0.1.0
// @description HTTP Micro service to make screenshot of a web page
// @termsOfService https://opensource.org/licenses/MIT
// @contact.name Bartosz Lenart
// @contact.email lenart.consulting@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host https://app-b33c1c94-0688-4054-92fd-c34a56577870.cleverapps.io
// @BasePath /
func main() {
	if isDebug {
		runtime.SetBlockProfileRate(1)
		debugger.Run()
	}
	router := mux.NewRouter()
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler).Methods("GET")
	screenshotRouter := router.PathPrefix("/screenshot").Subrouter()
	screenshotRouter.Use(validateToken)
	screenshotRouter.HandleFunc("/jpeg", screenshotJpeg).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         srvAddressAndPort,
		WriteTimeout: time.Second * srvReadWriteTimeout,
		ReadTimeout:  time.Second * srvReadWriteTimeout,
		IdleTimeout:  time.Second * srvIdleTimeout,
	}

	go func() {
		log.Printf("server started on %s\n", srvAddressAndPort)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("error creating server: %s\n", err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	close(c)
	wait := time.Duration(exitTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("server shut down")
	os.Exit(0)
}

// screenshotJpeg godoc
// @Summary Makes web page screenshot to jpeg
// @Description Makes full page screenshot to jpeg and returns jpeg buffer
// @Tags Scanners
// @Accept jpeg
// @Param token query string true "Token"
// @Param schema body URLRequestBody true "URL schema to screenshot a web page from"
// @Success 200 {} OK
// @Failure 401 {} Not authorized
// @Failure 500 {} Internal server error, if not valid query provided
// @Router /screenshot/jpeg [post]
func screenshotJpeg(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("error in scanPng decoding body json request: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, pageScreenshot(request.URL, 90, &buf)); err != nil {
		log.Printf("error in screenshot URL: %s, error: %s\n", request.URL, err)
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	if _, err := w.Write(buf); err != nil {
		log.Printf("error in screenshot when sending image screenshot buffer, error: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// validateToken performs token validation provided in URL like so ?token=this_micro_services_token
func validateToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request from IP: %s\n", r.RemoteAddr)
		token := r.URL.Query().Get("token")

		if token == "" || token != appToken {
			log.Printf("token not valid: %s\n", token)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// pageScreenshot takes a screenshot of the entire browser viewport.
// Note: this will override the viewport emulation settings.
func pageScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
