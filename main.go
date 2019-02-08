package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/cjburchell/tools-go/env"

	"github.com/cjburchell/go-uatu"

	"github.com/cjburchell/restmock/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	configFile := env.Get("CONFIG_FILE", "config.json")

	err := config.Setup(configFile)
	if err != nil {
		log.Fatal(err, "Unable to setup config File")
		return
	}

	port := env.Get("PORT", "8080")
	startHTTPTestEndpoints(port)
}

func startHTTPTestEndpoints(port string) {
	endpoints, err := config.GetEndpoints()
	if err != nil {
		log.Fatal(err, "Unable build endpoints")
		return
	}

	r := mux.NewRouter()

	for _, endpoint := range endpoints {
		r.HandleFunc(endpoint.Path, func(w http.ResponseWriter, r *http.Request) {
			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Error(err, "Unable to dump Request")
			}

			log.Print(string(requestDump))

			w.WriteHeader(endpoint.Response)
			w.Header().Set("Content-Type", endpoint.ContentType)
			_, err = w.Write([]byte(endpoint.ResponseBody))
			if err != nil {
				log.Error(err, "Unable to write response")
			}

		}).Methods(endpoint.Method)
	}

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Started http Server on port: " + port)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
