package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var imagePath, imageExt string

func main() {
	// Register static file server for images once
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	// Register main handler
	http.HandleFunc("/", handleRequest)

	// Start server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Server failed to start: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("Handling request for path: %s", path)

	// Handle root path
	if path == "/" || path == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Println("Serving JSON for root path")
		json.NewEncoder(w).Encode(map[string]string{"message": "the fuck r u doing here?"})
		return
	}

	// Extract status code from path (e.g., /404 -> 404)
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 1 {
		log.Printf("Invalid path format: %s", path)
		sendJSONError(w, http.StatusBadRequest, "Invalid path format")
		return
	}

	statusCode, err := strconv.Atoi(parts[0])
	if err != nil || statusCode < 100 || statusCode > 599 {
		log.Printf("Invalid status code: %s", parts[0])
		sendJSONError(w, http.StatusBadRequest, "Invalid status code")
		return
	}

	// Check if client accepts HTML
	acceptsHTML := strings.Contains(r.Header.Get("Accept"), "text/html")
	if acceptsHTML {
		log.Printf("Serving HTML for status code: %d", statusCode)
		serveHTML(w, statusCode)
	} else {
		log.Printf("Serving JSON for status code: %d", statusCode)
		sendJSONError(w, statusCode, http.StatusText(statusCode))
	}
}

func sendJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": statusCode,
		"error":  message,
	})
}

func serveHTML(w http.ResponseWriter, statusCode int) {
	// Check for image with .png, .jpg, or .gif extension
	extensions := []string{".png", ".jpg", ".gif"}
	imageExists := false
	for _, ext := range extensions {
		path := filepath.Join("images", fmt.Sprintf("%d%s", statusCode, ext))
		log.Printf("Checking image existence at path: %s", path)
		if _, err := os.Stat(path); err == nil {
			imagePath = path
			imageExt = ext
			imageExists = true
			log.Println("Image exists")
			break
		} else {
			log.Printf("Image check failed for %s: %v", path, err)
		}
	}

	// Load and parse the HTML template
	t, err := template.ParseFiles("template.html")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Failed to render template")
		return
	}

	// Render HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	data := struct {
		StatusCode   int
		ErrorMessage string
		ImageExists  bool
		ImageExt     string
	}{
		StatusCode:   statusCode,
		ErrorMessage: http.StatusText(statusCode),
		ImageExists:  imageExists,
		ImageExt:     imageExt,
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Failed to execute template")
	}
}
