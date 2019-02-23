package mock

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"context"

	"github.com/cjburchell/go-uatu"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	r         *mux.Router
	srv       *http.Server
	endpoints []*Endpoint
}

func (s *Server) Start(port int) {
	if s.r == nil {
		s.r = mux.NewRouter()
		log.Warn("Starting Server with no endpoints")
	}

	loggedRouter := handlers.LoggingHandler(os.Stdout, s.r)

	s.srv = &http.Server{
		Handler:      loggedRouter,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Started http Server on port: %d", port)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func (s *Server) Endpoint(name, method, path string) *Endpoint {
	endpoint := Endpoint{name: name, path: path, method: method}

	if s.r == nil {
		s.r = mux.NewRouter()
	}

	if len(endpoint.path) == 0 || len(endpoint.method) == 0 {
		return nil
	}

	endpoint.Reply = &Reply{}
	endpoint.Handler = endpoint.Reply

	log.Printf("Loading Endpoint %s %s %s", endpoint.name, endpoint.method, endpoint.path)
	endpoint.route = s.r.HandleFunc(endpoint.path, handleEndpoint(&endpoint)).Methods(endpoint.method)

	if s.endpoints == nil {
		s.endpoints = make([]*Endpoint, 0)
	}

	s.endpoints = append(s.endpoints, &endpoint)

	return &endpoint
}

func handleEndpoint(endpoint *Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling Endpoint %s %s %s", endpoint.name, endpoint.method, endpoint.path)

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

		endpoint.Handler.Handle(w, r)
	}
}
