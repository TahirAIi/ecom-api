format:
	go fmt -w .

generate-docs:
	GOROOT=/usr/local/go swagger generate spec -w ./cmd/api/ -o ./swagger.yaml --scan-models

serve-docs:
	swagger serve swagger.yaml --flavor=swagger --port=8080 --no-open

run:
	format generate-docs serve-docs