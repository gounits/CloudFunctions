package main

import (
	"log"
	"net/http"

	"github.com/gounits/CloudFunctions/api"
)

func main() {
	http.HandleFunc("POST /api/translate", api.Translate)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
