# Lacon HTTP API

Lacon provides an HTTP server for file compression and decompression using Huffman coding. This document describes the available API endpoints.

## Base URL

The server runs on `http://localhost:8080` by default. Set the `PORT` environment variable to change the port.

```
http://localhost:8080
```

## Common Headers

| Header         | Value                     |
|----------------|---------------------------|
| `Content-Type` | `multipart/form-data`     |
| `Origin`       | Any (CORS enabled)        |

The server enables CORS with `Access-Control-Allow-Origin: *`, making it accessible from web frontends running on any domain.

---

## POST /encode

Compresses a file using Huffman encoding and returns the compressed `.lacon` file.

### Request

- **Method:** `POST`
- **Content-Type:** `multipart/form-data`
- **Body parameter:** `file` — The file to compress (any text file)

### Responses

#### ✅ 200 OK — Successful compression

The response is a binary `.lacon` file download.

| Header                | Value                                                  |
|-----------------------|--------------------------------------------------------|
| `Content-Type`        | `application/octet-stream`                             |
| `Content-Disposition` | `attachment; filename="<original_filename>.lacon"`     |
| `Content-Length`      | Size of compressed data in bytes                       |

#### ❌ 400 Bad Request — Missing file or encoding failure

```json
{
  "error": "missing 'file' field in form data: ..."
}
```

```json
{
  "error": "encoding failed: ..."
}
```

#### ❌ 405 Method Not Allowed

```json
{
  "error": "only POST method is allowed"
}
```

### cURL Example

```bash
# Compress a file
curl -X POST http://localhost:8080/encode \
  -F "file=@document.txt" \
  --output document.txt.lacon
```

---

## POST /decode

Decompresses a `.lacon` file back to its original content.

### Request

- **Method:** `POST`
- **Content-Type:** `multipart/form-data`
- **Body parameter:** `file` — The `.lacon` compressed file to decompress

### Responses

#### ✅ 200 OK — Successful decompression

The response is a binary file download with the original content.

| Header                | Value                                                      |
|-----------------------|------------------------------------------------------------|
| `Content-Type`        | `application/octet-stream`                                 |
| `Content-Disposition` | `attachment; filename="<original_filename>.lacon.decoded"` |
| `Content-Length`      | Size of decompressed data in bytes                         |

#### ❌ 400 Bad Request — Missing file or decoding failure

```json
{
  "error": "missing 'file' field in form data: ..."
}
```

```json
{
  "error": "decoding failed: ..."
}
```

#### ❌ 405 Method Not Allowed

```json
{
  "error": "only POST method is allowed"
}
```

### cURL Example

```bash
# Decompress a .lacon file
curl -X POST http://localhost:8080/decode \
  -F "file=@document.txt.lacon" \
  --output restored.txt
```

---

## Full Workflow Example

```bash
# Start the server
go run ./server/main.go &

# Encode a file
echo "Hello, Lacon!" > hello.txt
curl -s -X POST http://localhost:8080/encode \
  -F "file=@hello.txt" \
  --output hello.txt.lacon

# Decode the compressed file
curl -s -X POST http://localhost:8080/decode \
  -F "file=@hello.txt.lacon" \
  --output hello_restored.txt

# Verify
diff hello.txt hello_restored.txt  # Should produce no output

# Stop the server
kill %1
```

---

## Error Codes Reference

| HTTP Status | Meaning                                |
|-------------|----------------------------------------|
| `200`       | Success — file processed and returned  |
| `400`       | Bad request — missing file field, invalid data, or processing failure |
| `405`       | Method not allowed — only `POST` is accepted |
| `500`       | Internal server error — unexpected failure |

---

## Postman Collection

A Postman collection is available at [`server/postman_collection.json`](../server/postman_collection.json). Import it into Postman to quickly test the API endpoints.

### Import Steps

1. Open Postman
2. Click **Import** → **Upload Files**
3. Select `server/postman_collection.json`
4. The collection includes pre-configured requests for both `/encode` and `/decode` with test scripts

---

## Running the Server

```bash
# From the project root
go run ./server/main.go

# Or build and run
go build -o lacon-server ./server
./lacon-server

# With a custom port
PORT=3000 go run ./server/main.go