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

// statusDescriptions maps HTTP status codes to their IANA-like descriptions
var statusDescriptions = map[int]string{
	100: "The server has received the request headers and the client should proceed to send the request body.",
	101: "The server is switching protocols as requested by the client.",
	102: "The server has received and is processing the request, but no response is available yet.",
	103: "The server is providing early hints about possible response headers.",
	200: "The request has succeeded.",
	201: "The request has been fulfilled and resulted in a new resource being created.",
	202: "The request has been accepted for processing, but the processing has not been completed.",
	203: "The response is a transformed version of the requested resource from another source.",
	204: "The server successfully processed the request, but is not returning any content.",
	205: "The server successfully processed the request, but the client should reset the document view.",
	206: "The server is delivering only part of the resource due to a range header sent by the client.",
	207: "The response provides status for multiple independent operations.",
	208: "The members of a DAV binding have been enumerated elsewhere in the response.",
	226: "The server has fulfilled a request for the resource, and the response is a representation of the result of one or more instance-manipulations.",
	300: "The request has multiple possible responses. User-agent or user should choose one of them.",
	301: "The requested resource has been permanently moved to a new URI.",
	302: "The requested resource resides temporarily under a different URI.",
	303: "The response to the request can be found under another URI using the GET method.",
	304: "The resource has not been modified since the last request.",
	305: "The requested resource must be accessed through the proxy given by the Location field.",
	306: "This status code is no longer used but reserved.",
	307: "The requested resource resides temporarily under a different URI, and the request should not be modified.",
	308: "The requested resource has been permanently moved to a new URI, and the request should not be modified.",
	400: "The server cannot or will not process the request due to an apparent client error.",
	401: "Authentication is required and has failed or has not yet been provided.",
	402: "Payment is required to access the requested resource.",
	403: "The server understood the request, but is refusing to fulfill it.",
	404: "The requested resource could not be found on the server.",
	405: "The request method is not supported for the requested resource.",
	406: "The requested resource is not available in a format acceptable to the client.",
	407: "The client must authenticate with a proxy before the request can proceed.",
	408: "The server timed out waiting for the client's request.",
	409: "The request could not be completed due to a conflict with the current state of the resource.",
	410: "The requested resource is permanently gone and will not be available again.",
	411: "The request requires a Content-Length header that was not provided.",
	412: "One or more preconditions in the request header fields evaluated to false.",
	413: "The request entity is larger than the server is willing or able to process.",
	414: "The request URI is too long for the server to process.",
	415: "The request's media type is not supported by the server or resource.",
	416: "The requested range is not satisfiable by the server.",
	417: "The expectation given in the Expect request header could not be met.",
	418: "The server refuses to brew coffee because it is a teapot.",
	421: "The request was directed at a server that is not able to produce a response.",
	422: "The request was well-formed but unable to be followed due to semantic errors.",
	423: "The resource that is being accessed is locked.",
	424: "The request failed because it depended on another request that failed.",
	425: "The server is unwilling to risk processing a request that might be replayed.",
	426: "The client must upgrade to a different protocol to continue.",
	428: "The server requires the request to be conditional.",
	429: "The client has sent too many requests in a given amount of time.",
	431: "The server is unwilling to process the request because its header fields are too large.",
	451: "The resource is unavailable due to legal reasons.",
	500: "The server encountered an unexpected condition that prevented it from fulfilling the request.",
	501: "The server does not support the functionality required to fulfill the request.",
	502: "The server received an invalid response from an upstream server.",
	503: "The server is currently unavailable due to maintenance or overloading.",
	504: "The server, acting as a gateway, did not receive a timely response from an upstream server.",
	505: "The server does not support the HTTP protocol version used in the request.",
	506: "The server has an internal configuration error where variant negotiation is self-referential.",
	507: "The server is unable to store the representation needed to complete the request.",
	508: "The server detected an infinite loop while processing the request.",
	510: "Further extensions to the request are required for the server to fulfill it.",
	511: "The client needs to authenticate to gain network access.",
}

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

	// Get the detailed description, falling back to http.StatusText if not defined
	description, exists := statusDescriptions[statusCode]
	if !exists {
		description = http.StatusText(statusCode)
	}

	// Render HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	data := struct {
		StatusCode   int
		ErrorMessage string
		Description  string
		ImageExists  bool
		ImageExt     string
	}{
		StatusCode:   statusCode,
		ErrorMessage: http.StatusText(statusCode),
		Description:  description,
		ImageExists:  imageExists,
		ImageExt:     imageExt,
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Failed to execute template")
	}
}
