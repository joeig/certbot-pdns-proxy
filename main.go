package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Error struct {
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func runServer() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/authenticate", Authenticate).Methods("POST")
	router.HandleFunc("/v1/cleanup", Cleanup).Methods("DELETE")
	router.HandleFunc("/v1/health", Health).Methods("GET")
	router.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServeTLS(C.Server.ListenAddress, C.Server.CertFile, C.Server.KeyFile, router))
}

func main() {
	configFile := flag.String("config", "config.yml", "Configuration file")
	flag.Parse()
	parseConfig(&C, configFile)
	runServer()
}
