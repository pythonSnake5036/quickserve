package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
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
	singleFilePtr := flag.String("file", "", "single file mode")
	flag.Parse()

	isHttps := *certFilePtr != "" && *keyFilePtr != ""
	isSingleFile := false
	singleFileUrl := "/" + *singleFilePtr

	if *singleFilePtr != "" {
		if _, err := os.Stat(*singleFilePtr); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("File %s does not exist\n", *singleFilePtr)
			os.Exit(1)
		} else {
			fmt.Printf("Using single file mode, redirecting all requests to %s\n", *singleFilePtr)
			isSingleFile = true
		}
	}

	if !isSingleFile {
		fmt.Printf("Serving %s\n", *dirPtr)
	}

	if isHttps {
		fmt.Printf("Listening on: %s:%d (https)\n", *hostPtr, *portPtr)
	} else {
		fmt.Printf("Listening on: %s:%d\n", *hostPtr, *portPtr)
	}

	fileServer := http.FileServer(http.Dir(*dirPtr))

	handler := func(w http.ResponseWriter, r *http.Request) {
		if isSingleFile {
			if r.URL.Path != singleFileUrl {
				w.Header().Set("Location", singleFileUrl)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}

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
