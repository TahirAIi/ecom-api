# How to run
If you have `make`:
- Run `make generate-docs` to generate swagger docs.
- Run `make serve-docs` to start docs server.
- Run `make serve-api` to start the api server.

If you don't have `make`:
- Run `GOROOT=/usr/local/go swagger generate spec -w ./cmd/api/ -o ./swagger.yaml --scan-models` to generate swagger docs.
- Run `swagger serve swagger.yaml --port=9000 --no-open`.
- Run `go run ./cmd/api` at the root of project to start api server.
