FROM golang:1.19-alpine as builder

RUN apk update && apk add --no-cache git=2.40.1-r0

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -a -o /app/bin/json-csv-converter . \
    && chmod +x /app/bin/json-csv-converter \
    && cp /app/bin/json-csv-converter /usr/local/bin/json-csv-converter

FROM gcr.io/distroless/static:nonroot
ENV LOCALES_PATH=/locales

WORKDIR /
COPY --from=builder /app/bin/json-csv-converter .
COPY ./translations/locales ./locales
USER 65532:65532

EXPOSE 8000
ENTRYPOINT ["/json-csv-converter"]
