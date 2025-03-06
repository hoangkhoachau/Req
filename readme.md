# HTTP Request CLI Tool

This is a command-line interface (CLI) tool written in Go that allows users to send HTTP requests (GET, POST, PUT, DELETE) to a specified URL with customizable headers, data, and output options. It provides a colorized output for better readability and supports JSON formatting for requests and responses. The tool is designed to be flexible, handling both simple requests and more complex scenarios with custom headers and data inputs.

## Features

- **Supported HTTP Methods**: GET, POST, PUT, DELETE
- **Custom Headers**: Add headers via command-line arguments (e.g., `key:value`).
- **Data Input**: Send data via stdin, command-line flags (`-d`), or file (`@filename`).
- **JSON Support**: Automatically detects and pretty-prints JSON responses; supports JSON data input.
- **Colorized Output**: Responses and requests are color-coded based on status codes and method types.
- **Output Options**: Save responses to a file (`-o`), print headers only (`-h`), or display full request/response details (`-f`).
- **Query Parameters**: Append query parameters using `key==value`.
- **Flexible URL Handling**: Automatically prepends `http://` if no protocol is specified; supports localhost shorthand (e.g., `:8080`).

## Installation

1. Ensure you have Go installed (version 1.16 or later recommended).
2. Clone or download this code into a file (e.g., `httpcli.go`).
3. Build the binary:
   ```bash
   go build -o httpcli httpcli.go
   ```
4. Move the binary to a directory in your PATH (optional):
   ```bash
   mv httpcli /usr/local/bin/
   ```

## Usage

### Basic Syntax
```
httpcli [METHOD] URL [OPTIONS] [HEADERS] [DATA]
```

- **`METHOD`**: The HTTP method (e.g., `GET`, `POST`, `PUT`, `DELETE`). Optional; defaults to `GET` if omitted or if data is not provided. If data is provided without a method, it defaults to `POST`.
- **`URL`**: The target URL. Required. Can be a full URL (e.g., `http://example.com`), a domain (e.g., `example.com`), or a localhost shorthand (e.g., `:8080/path`).
- **`OPTIONS`**: Command-line flags to modify behavior (see below).
- **`HEADERS`**: Custom headers in `key:value` format.
- **`DATA`**: Data to send in the request body or as query parameters (see formats below).

### Options
- **`-h`**: Prints only the response headers, omitting the body.
- **`-f`**: Prints full details of both the request and response, including headers and body.
- **`-d <data>`**: Specifies the request body data inline. Must be followed by a string (e.g., `-d "key=value"` or `-d '{"key": "value"}'`).
- **`-o <filename>`**: Saves the response body to the specified file. Must be followed by a valid filename.
- **`@filename`**: Reads the request body from the specified file (e.g., `@data.json`).

### Positional Arguments
The tool interprets positional arguments based on their format:
- **`key:value`**: Adds a custom header (e.g., `Authorization:Bearer-token`).
- **`key==value`**: Appends a query parameter to the URL (e.g., `id==123` becomes `?id=123`).
- **`key=value`**: Adds a key-value pair to the request body (non-JSON format).
- **`key:=value`**: Adds a JSON object to the request body, where `value` must be valid JSON (e.g., `user:={"name": "John"}`).

### URL Handling
- If the URL doesn’t include a protocol (e.g., `http://` or `https://`), `http://` is prepended automatically.
- If the URL starts with `:`, it’s treated as a localhost shorthand:
  - `:8080` becomes `http://localhost:8080`.
  - `:8080/path` becomes `http://localhost:8080/path`.

### Data Input Methods
You can provide data in multiple ways, but only one method can be used per request:
1. **Inline with `-d`**:
   ```bash
   httpcli POST http://example.com -d '{"name": "John"}'
   ```
2. **From a File with `@filename`**:
   ```bash
   httpcli POST http://example.com @data.json
   ```
3. **Via Stdin**:
   ```bash
   echo '{"key": "value"}' | httpcli POST http://example.com
   ```
4. **Using Positional Arguments**:
   ```bash
   httpcli POST http://example.com name=John age=30
   ```

### Detailed Examples

#### 1. Basic GET Request
```bash
httpcli GET http://example.com
```
- Sends a GET request to `http://example.com`.
- Prints the response body (or headers if `-h` is used).

#### 2. GET with Query Parameters
```bash
httpcli http://example.com id==123 name==John
```
- Sends a GET request to `http://example.com?id=123&name=John`.
- Omitting `GET` is valid; it defaults to GET when no data is provided.

#### 3. POST with Inline JSON Data
```bash
httpcli POST http://example.com -d '{"name": "John", "age": 30}'
```
- Sends a POST request with a JSON body.
- The `Content-Type` header is automatically set to `application/json`.

#### 4. POST with Data from a File
```bash
httpcli POST http://example.com @data.json
```
- Reads `data.json` and sends its contents as the request body.
- If the file contains valid JSON, `Content-Type` is set to `application/json`.

#### 5. POST with Headers
```bash
httpcli POST http://example.com Authorization:Bearer-token Content-Type:application/json -d '{"key": "value"}'
```
- Adds custom headers (`Authorization` and `Content-Type`) to the request.
- Sends the JSON data provided with `-d`.

#### 6. PUT with Key-Value Data
```bash
httpcli PUT http://example.com name=John age=30
```
- Constructs a request body with `name=John` and `age=30`.
- Since it’s not JSON, `Content-Type` is determined by `http.DetectContentType`.

#### 7. DELETE with Full Output
```bash
httpcli DELETE http://example.com -f
```
- Sends a DELETE request and prints both the request and response details.
- Useful for debugging.

#### 8. Save Response to File
```bash
httpcli GET http://example.com -o response.txt
```
- Saves the response body to `response.txt`.
- No output is printed to the terminal unless `-f` is also used.

#### 9. Headers Only
```bash
httpcli GET http://example.com -h
```
- Prints only the response headers, omitting the body.

#### 10. Piping JSON via Stdin
```bash
echo '{"id": 1, "title": "Test"}' | httpcli POST http://example.com
```
- Sends a POST request with the piped JSON data.
- Automatically sets `Content-Type` to `application/json`.

#### 11. Complex JSON with `:=` Syntax
```bash
httpcli POST http://example.com user:={"name": "John", "age": 30}
```
- Constructs a JSON body: `{"user": {"name": "John", "age": 30}}`.
- Defaults to POST since data is provided.

#### 12. Localhost Shorthand
```bash
httpcli :8080/api endpoint==test
```
- Sends a GET request to `http://localhost:8080/api?endpoint=test`.

#### 13. File-Based Query Parameter
```bash
httpcli http://example.com id==@id.txt
```
- Reads the contents of `id.txt` and appends it as a query parameter (e.g., `?id=<contents>`).

### Edge Cases
- **Multiple Data Sources**: Combining `-d`, `@filename`, and stdin will result in an error ("Invalid data"). Choose one method.
- **Invalid JSON**: If JSON data is malformed, an error is printed, and the program exits.
- **Missing URL**: Omitting the URL results in "Missing url" error.
- **Unsupported Method**: Only GET, POST, PUT, and DELETE are supported; others return "unsupported method".

### Tips
- Use `-f` to debug requests and see exactly what’s being sent and received.
- Quote data strings with spaces (e.g., `-d "key=value with spaces"`).
- Ensure files referenced with `@filename` exist and are readable.

## Output Formatting

- **Request Output** (with `-f`):
  - Method: Color-coded (e.g., GET=Green, POST=Blue, DELETE=Red).
  - URL and protocol.
  - Headers in green.
  - Body (pretty-printed if JSON).
- **Response Output**:
  - Status code: Color-coded (2xx=Green, 3xx=Blue, 4xx=Yellow, 5xx=Red).
  - Headers in green (if `-f` or error).
  - Body: Pretty-printed JSON or raw text (unless `-h`).

## Color Codes
- **Methods**: GET (Green), POST (Blue), PUT (Yellow), DELETE (Red).
- **Status Codes**: 2xx (Green), 3xx (Blue), 4xx (Yellow), 5xx (Red).
- **JSON Values**: Strings (Green), Numbers (Yellow), Booleans (Blue), Null (Gray).

## Code Structure

### Main Components
- **`sendRequest`**: Constructs and sends the HTTP request with headers and data.
- **`printRequest`**: Displays the request details (with `-f`).
- **`printRespond`**: Prints the response, handling status codes, headers, and body.
- **`prettyJson`**: Formats JSON data with indentation and colorized values.
- **`main`**: Parses arguments, processes flags, and orchestrates the request.

### Dependencies
- Standard Go libraries only (`net/http`, `encoding/json`, etc.).

## Limitations

- Only supports GET, POST, PUT, and DELETE methods.
- No support for PATCH (despite color coding in `printRequest`).
- Limited error handling for malformed JSON or invalid files.
- No timeout configuration for requests.

## Future Improvements

- Add support for PATCH and other HTTP methods.
- Implement request timeout options.
- Enhance error messages for better debugging.
- Add support for multipart/form-data requests.
- Include a `--help` flag with detailed usage info.

## License

This project is unlicensed and provided as-is for educational purposes. Feel free to modify and distribute it as needed.
