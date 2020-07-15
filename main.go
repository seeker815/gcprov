package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/seeker815/gcprov/gc"
)

func main() {
	fmt.Println("gcProv service handles instance healtcheck, allows to provision instances on Google Cloud")
	router := gc.Router()
	log.Fatal(http.ListenAndServe(":8000", router)
}