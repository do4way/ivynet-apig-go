package main

import (
	"log"
	"net/http"
	"os"

	"github.com/do4way/ivynet-apig-go/eimbc"
	"github.com/gorilla/mux"
)

var port = os.Getenv("HTTP_PORT")

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{topic}/{id}", eimbc.HTTPPostHandler).Methods("POST")
	router.HandleFunc("/{topic}/{id}", eimbc.HTTPGetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
