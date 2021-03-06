# GO Basilisk

![Build](https://github.com/bartOssh/go_basilisk/workflows/Go/badge.svg?branch=main)

HTTP micro-service to make screenshot of a web page to jpeg image 

[![FUN](https://varsisava.pl/wp-content/uploads/2016/12/Operacja-Bazyliszek.jpg)](https://www.youtube.com/watch?v=qS2xTGLCu-M&t)


## Generate swagger docs

Run before first build or run to generate docs (swagger)

```bash
  go get -u github.com/swaggo/swag/cmd/swag
  swag init -g main.go
```

## Development and deployment

1. Project is developed in GO, run:
   ```bash
        go run .
   ```
2. Build with:
    ```bash
        go build -o <name> .
    ```
3. Destination machine is Ubuntu 20.04^, but should work on other Linux OS, Windows and Darwin
4. Project aims to be deployed as binary executable
5. Please set environment variables according to .env.example or make .env alongside executable

## Tests

1. Run `go test`
2. Tests are covering one existing endpoint `/screenshot/jpeg?token=your_app_token` and services/helpers units

## Issues, pull requests and suggestion

Just open pull request with proposed changes or add [issue](https://github.com/bartOssh/go_basilisk/issues)

## License

[MIT](https://opensource.org/licenses/MIT)

## Documentation

Run project open: `http://localhost:8888/docs/index.html`, to view API documentation

## Performance

Performance highly depands on how fast web page We want to take screenshot of is responding to request.
Making screenshot of: "https://github.com/bartOssh/go_basilisk" took ~1.6 second on average.
Tested on shared 1 CPU server with 512 MB RAM.
 
