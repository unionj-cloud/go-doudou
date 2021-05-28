package httpsrv

import "net/http"

type Interface interface {
	// Run the service
	Run()
	// Register routes
	Route(route Route)
	// Create http server
	NewServer(router http.Handler) *http.Server
	// Use middleware
	Use(mwf ...func(http.Handler) http.Handler)
}
