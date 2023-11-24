build:
	@docker build -t json-csv-converter:latest .

install:
	@go get ./...

run:
	@go run .

run-docker:
	@docker run -p 8000:8000 json-csv-converter

build-and-run-docker:
	@docker build -t json-csv-converter:latest .
	@docker run -p 8000:8000 json-csv-converter

test:
	@go test ./...
