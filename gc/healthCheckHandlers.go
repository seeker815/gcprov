package gc

import (
	"encoding/json"
	"log"
	"net/http"
)

// GetHealthCheck implements a basic health check
func GetHealthCheck(writer http.ResponseWriter, request *http.Request) {
	InitHeaders(writer)
	json.NewEncoder(writer).Encode("HealthCheck ran successfully!")
	log.Print("Healtcheck invoked")
}

// InitHeaders initializes headers
func InitHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}
