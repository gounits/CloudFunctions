package main

import (
	"github.com/gounits/CloudFunctions/tool/middleware"
	"log"
	"net/http"

	"github.com/gounits/CloudFunctions/api"
)

func main() {
	middlewares := middleware.LoggingMiddleware(http.DefaultServeMux.ServeHTTP)
	http.HandleFunc("POST /api/translate/{name...}", api.Translate)
	log.Fatal(http.ListenAndServe(":8080", middlewares))
}
