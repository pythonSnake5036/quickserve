package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	portFlag int
	dirFlag  string
)

var compressionTypes = map[string]string{
	".gz": "gzip",
}

var contentTypes = map[string]string{
	".wasm": "application/wasm",
}

func applyHeaders(w http.ResponseWriter, path string) {
	for ext, compressionType := range compressionTypes {
		if strings.HasSuffix(path, ext) {
			w.Header().Set("Content-Encoding", compressionType)

			base := strings.TrimSuffix(path, ext)
			applyHeaders(w, base)
		}
	}

	for ext, contentType := range contentTypes {
		if strings.HasSuffix(path, ext) {
			w.Header().Set("Content-Type", contentType)
		}
	}
}

func main() {
	flag.IntVar(&portFlag, "port", 8080, "port to listen")
	flag.StringVar(&dirFlag, "dir", ".", "directory to serve")
	flag.Parse()

	fmt.Printf("Serving %s\nListening on port: %d\n", dirFlag, portFlag)

	fileServer := http.FileServer(http.Dir(dirFlag))

	handler := func(w http.ResponseWriter, r *http.Request) {
		applyHeaders(w, r.URL.Path)

		fileServer.ServeHTTP(w, r)
	}

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portFlag), nil))
}
