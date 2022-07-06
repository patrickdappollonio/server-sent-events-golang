package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://localhost:8080/logs")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		panic(err)
	}
}
