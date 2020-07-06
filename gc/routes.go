package gc

import (
	"github.com/gorilla/mux"
)

// Router implements the routes for Provisioner service
func Router() *mux.Router {
	router := mux.NewRouter()
	prefix := "/v1/instances"

	router.HandleFunc("/healthcheck", GetHealthCheck).Methods("GET")
	router.HandleFunc(prefix+"/create", CreateInstance).Queries(
		"project", "{project}",
		"zone", "{zone}",
		"username", "{username}",
		"userpass", "{userpass}",
	).Methods("POST")
	router.HandleFunc(prefix+"/status", GetInstanceStatus).Queries(
		"project", "{project}",
		"zone", "{zone}",
		"instance", "{instance}",
	).Methods("GET")

	return router
}
