# GO Basilisk v 0.1.0

HTTP micro-service scanning web page to png format

[![FUN](https://varsisava.pl/wp-content/uploads/2016/12/Operacja-Bazyliszek.jpg)](https://www.youtube.com/watch?v=JT9e4qzLdrM)

## If swagger not generated

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
2. Destination machine is Ubuntu 20.04^, but should work on other Linux OS, Windows and Darwin
3. Project is deployed as binary executable
4. Please set environment variables according to .env.example or set global ENV variables alongside executable

## Test

1. Run `go test .`

## Issues, pull requests and suggestion

Just open pull request with proposed changes or add [issue](https://github.com/bartOssh/go_basilisk/issues)

## License

[MIT](https://opensource.org/licenses/MIT)

## Documentation

Run project open: `http://localhost:8888/docs/index.html`, to view API documentation

## Why

I am a big fun(c) of the Youtube podcast series: [Just for func](https://www.youtube.com/channel/UC_BzFbxG2za3bp5NRRRXJSw), and I enjoy GO.
