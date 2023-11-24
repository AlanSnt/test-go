# JSON-CSV Converter

This microservice has been developed in Go to facilitate data export in CSV or Excel formats. It can be used to generate CSV files from structured data and download them as attachments in HTTP responses.

## Context of Use

This microservice is useful when you need to generate CSV files from structured data in your application and provide them for download to users. It is particularly suited for web applications that require data export for reports or further analysis.

## Commands

To run this microservice, you need to have Go installed on your machine. Here are the steps to run it:

1. Clone this repository:

   ```shell
   git clone git@git.pw.fr:unyc/code/services/unyc/json-csv-converter/json-csv-converter.git
   cd json-csv-converter
   ```

2. Run the microservice:

   ```shell
   go run main.go
   ```

The microservice will run on port `8000` by default. You can customize it by modifying the line `router.Run(":8000")` in the `main.go` file.

## Endpoints

The microservice exposes a single endpoint for data export:

- `POST /export?type={format}`

  This endpoint allows you to specify the output format, which can be "csv" or "excel." You need to send a POST request with a JSON body containing the data to export in the following format:

  ```json
  {
    "fileName": "file_name",
    "records": [ ... ],
    "columns": [ ... ],
    "delimiter": "," // Optional default value is ;
  }
  ```

  - `"fileName"`: The name of the output file (without the extension).
  - `"records"`: An array of records (each record being a JSON object).
  - `"columns"`: An array of column names corresponding to the data.

## Example Request

Here's an example HTTP request to export data in CSV format:

```shell
curl -X POST http://localhost:8000/export?type=csv -H "Content-Type: application/json" -d '{
  "fileName": "data_export",
  "records": [
    {
      "name": "John Doe",
      "age": 30,
      "email": "john@example.com"
    },
    {
      "name": "Jane Smith",
      "age": 25,
      "email": "jane@example.com"
    }
  ],
  "columns": ["name", "age", "email"]
}'
```

## Supported Output Formats

- CSV: The data will be generated in CSV format and returned as a downloadable CSV attachment.

## Possible Errors

- If the specified format is not supported, you will receive a response with HTTP status code `415 Unsupported Media Type`.
- If the specified format is "excel," you will receive a response with HTTP status code `501 Not Implemented`.

## Dependencies

This microservice uses the Gin and Cors libraries to handle routes and CORS permissions. Make sure to install them before running the microservice.
