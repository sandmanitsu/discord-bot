package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleOAuth2Callback)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	fmt.Fprintf(w, "Authorization code: %s", code)
}
