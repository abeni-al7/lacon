package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/abeni-al7/lacon/core"
)

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "only POST method is allowed")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse multipart form: %v", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing 'file' field in form data: %v", err)
		return
	}
	defer file.Close()

	inputData, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read uploaded file: %v", err)
		return
	}

	var outputBuf bytes.Buffer
	if err := core.Encode(bytes.NewReader(inputData), &outputBuf); err != nil {
		writeError(w, http.StatusInternalServerError, "encoding failed: %v", err)
		return
	}

	filename := header.Filename
	if filename == "" {
		filename = "output"
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.lacon"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", outputBuf.Len()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, &outputBuf)
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "only POST method is allowed")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse multipart form: %v", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing 'file' field in form data: %v", err)
		return
	}
	defer file.Close()

	inputData, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read uploaded file: %v", err)
		return
	}

	var outputBuf bytes.Buffer
	if err := core.Decode(bytes.NewReader(inputData), &outputBuf); err != nil {
		writeError(w, http.StatusInternalServerError, "decoding failed: %v", err)
		return
	}

	filename := header.Filename
	if filename == "" {
		filename = "output"
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.decoded"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", outputBuf.Len()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, &outputBuf)
}

func writeError(w http.ResponseWriter, statusCode int, format string, args ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": fmt.Sprintf(format, args...),
	})
}