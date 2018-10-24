package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/cjburchell/restmock/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	config.Setup("config.json")

	go startHTTPTestEndpoints()

	startHTTPConfigEndpoints()
}

func handleInfo(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func startHTTPConfigEndpoints() {
	r := mux.NewRouter()
	r.HandleFunc("/info", handleInfo).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8082",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf(err.Error())
	}
}

func startHTTPTestEndpoints() *http.Server {
	endpoints, _ := config.GetEndpoints()

	r := mux.NewRouter()

	for _, endpoint := range endpoints {
		r.HandleFunc(endpoint.Path, func(w http.ResponseWriter, r *http.Request) {
			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(requestDump))

			w.WriteHeader(endpoint.Response)
			w.Header().Set("Content-Type", endpoint.ContentType)
			w.Write([]byte(endpoint.ResponseBody))
		}).Methods(endpoint.Method)
	}

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         ":8088",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Started http Server")

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf(err.Error())
	}

	return srv
}
