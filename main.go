package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	portPtr := flag.Int("port", 8080, "port to listen")
	dirPtr := flag.String("dir", ".", "directory to serve")
	certFilePtr := flag.String("cert", "", "certificate file to use for TLS")
	keyFilePtr := flag.String("key", "", "key file to use for TLS")
	flag.Parse()

	isHttps := *certFilePtr != "" && *keyFilePtr != ""

	if isHttps {
		fmt.Printf("Serving %s\nListening on port: %d (https)\n", *dirPtr, *portPtr)
	} else {
		fmt.Printf("Serving %s\nListening on port: %d\n", *dirPtr, *portPtr)
	}

	fileServer := http.FileServer(http.Dir(*dirPtr))

	handler := func(w http.ResponseWriter, r *http.Request) {
		applyHeaders(w, r.URL.Path)

		fileServer.ServeHTTP(w, r)
	}

	http.HandleFunc("/", handler)

	if isHttps {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *portPtr), *certFilePtr, *keyFilePtr, nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), nil))
	}
}
