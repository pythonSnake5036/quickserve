package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var compressionTypes = map[string]string{
	".gz": "gzip",
}

func applyHeaders(w http.ResponseWriter, path string) {
	for ext, compressionType := range compressionTypes {
		if strings.HasSuffix(path, ext) {
			w.Header().Set("Content-Encoding", compressionType)

			base := strings.TrimSuffix(path, ext)
			applyHeaders(w, base)
		}
	}

	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}
}

func main() {
	hostPtr := flag.String("host", "127.0.0.1", "host to listen on")
	portPtr := flag.Int("port", 8080, "port to listen")
	dirPtr := flag.String("dir", ".", "directory to serve")
	certFilePtr := flag.String("cert", "", "certificate file to use for TLS")
	keyFilePtr := flag.String("key", "", "key file to use for TLS")
	flag.Parse()

	isHttps := *certFilePtr != "" && *keyFilePtr != ""

	if isHttps {
		fmt.Printf("Serving %s\nListening on: %s:%d (https)\n", *dirPtr, *hostPtr, *portPtr)
	} else {
		fmt.Printf("Serving %s\nListening on: %s:%d\n", *dirPtr, *hostPtr, *portPtr)
	}

	fileServer := http.FileServer(http.Dir(*dirPtr))

	handler := func(w http.ResponseWriter, r *http.Request) {
		applyHeaders(w, r.URL.Path)

		fileServer.ServeHTTP(w, r)
	}

	http.HandleFunc("/", handler)

	if isHttps {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *hostPtr, *portPtr), *certFilePtr, *keyFilePtr, nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *hostPtr, *portPtr), nil))
	}
}
