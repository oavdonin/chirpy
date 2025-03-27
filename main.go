package main

import (
	"net/http"
)

type webServer struct{}

func (webServer) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", webServer{})
	http.ListenAndServe(":8080", mux)
}
