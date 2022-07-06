package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, supported := w.(http.Flusher)

		if !supported {
			http.Error(w, "Streaming logs not supported by the client", http.StatusPreconditionFailed)
			return
		}

		scanner := bufio.NewScanner(os.Stdin)
		byteCount := 0
		for scanner.Scan() {
			byteCount += len(scanner.Bytes())
			fmt.Fprintf(w, "stdin -> %s\n", scanner.Text())
			flusher.Flush()
		}

		if err := scanner.Err(); err != nil {
			http.Error(w, "Can't read from stdin: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "info ->", byteCount, "bytes written to stdout")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Server listening on: %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
