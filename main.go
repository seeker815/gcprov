package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/seeker815/gcprovision/gc"
)

func main() {
	router := gc.Router()

	log.Fatal(http.ListenAndServe(":8000", router)
}