package main

import (
	"encoding/json"
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
		log.Printf("Loading Endpoint %s %s %s", endpoint.Name, endpoint.Method, endpoint.Path)
		r.HandleFunc(endpoint.Path, HandleEndpoint(endpoint)).Methods(endpoint.Method)
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

func HandleEndpoint(endpoint config.Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Error(err, "Unable to dump Request")
		}

		log.Print(string(requestDump))

		vars := mux.Vars(r)
		if len(vars) != 0 {
			log.Print("Values:")
			for key, value := range vars {
				log.Printf("Key: %s, Value: $s", key, value)
			}
		}

		if endpoint.Wait != 0 {
			log.Printf("Sleeping for %d Seconds", endpoint.Wait)
			time.Sleep(time.Duration(endpoint.Wait) * time.Second)
		}

		log.Printf("Send Response: %d %s Body: %s", endpoint.Response, endpoint.ContentType, endpoint.ResponseBody)
		w.WriteHeader(endpoint.Response)
		w.Header().Set("Content-Type", endpoint.ContentType)
		if endpoint.Header != nil {
			log.Print("Header")
			for key, value := range endpoint.Header {
				log.Printf("%s %s", endpoint.Response, endpoint.ContentType, endpoint.ResponseBody)
				w.Header().Set(key, value)
			}
		}

		if endpoint.ContentType == "application/json" {
			_, err = w.Write(endpoint.ResponseBody)
		} else {
			var body string
			err := json.Unmarshal(endpoint.ResponseBody, &body)
			if err == nil {
				_, err = w.Write([]byte(body))
			}
		}

		if err != nil {
			log.Error(err, "Unable to write response")
		}
	}
}
