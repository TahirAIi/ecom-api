generate-docs:
	GOROOT=/usr/local/go swagger generate spec -w ./cmd/api/ -o ./swagger.yaml --scan-models

serve-docs:
	swagger serve swagger.yaml --port=9000 --no-open

serve-api:
	go run ./cmd/api
