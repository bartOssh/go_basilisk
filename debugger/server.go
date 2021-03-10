package debugger

import (
	"fmt"
	"log"
	"net/http"
)

const (
	// DebugAddr is debugger server address
	DebugAddr = "127.0.0.1"
	// DebugPort is debugger server port
	DebugPort = 6060
)

type debugServer struct {
	*http.Server
}

func newDebugServer(address string) *debugServer {
	return &debugServer{
		&http.Server{
			Addr:    address,
			Handler: http.DefaultServeMux,
		},
	}
}

// Run runs debugger server concurrently
func Run() {
	debugServer := newDebugServer(fmt.Sprintf("%s:%d", DebugAddr, DebugPort))
	log.Printf("debugger started on %s:%d", DebugAddr, DebugPort)
	go func() {
		log.Fatal(debugServer.ListenAndServe())
	}()
}
