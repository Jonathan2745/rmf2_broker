package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		http.Error(w, `{"detail":"marshal failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"detail": msg})
}

// ServeGLB serves filePath, falling back to fallback if filePath is not found.
// Pass an empty fallback to return 404 when the file is missing.
func ServeGLB(w http.ResponseWriter, filePath, fallback string) {
	f, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && fallback != "" {
			f, err = os.Open(fallback)
		}
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				WriteError(w, http.StatusNotFound, "file not found")
			} else {
				WriteError(w, http.StatusInternalServerError, "failed to open file")
			}
			return
		}
	}
	defer f.Close()

	w.Header().Set("Content-Type", "model/gltf-binary")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}
