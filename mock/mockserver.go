package mock

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"context"

	"github.com/cjburchell/go-uatu"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type IServer interface {
	Start(port int)
	Stop() error
	Endpoint(name, method, path string) IEndpoint
	Verify() error
}

type server struct {
	r         *mux.Router
	srv       *http.Server
	endpoints []*endpoint
}

func CreateServer() IServer {
	return &server{}
}

func (s *server) Start(port int) {
	if s.r == nil {
		s.r = mux.NewRouter()
		log.Warn("Starting server with no endpoints")
	}

	loggedRouter := handlers.LoggingHandler(os.Stdout, s.r)

	s.srv = &http.Server{
		Handler:      loggedRouter,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Started http server on port: %d", port)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()
}

func (s *server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func (s *server) Endpoint(name, method, path string) IEndpoint {
	newEndpoint := createDefaultEndpoint(name, method, path)

	if s.r == nil {
		s.r = mux.NewRouter()
	}

	if len(newEndpoint.path) == 0 || len(newEndpoint.method) == 0 {
		return nil
	}

	log.Printf("Loading newEndpoint %s %s %s", newEndpoint.name, newEndpoint.method, newEndpoint.path)
	newEndpoint.route = s.r.HandleFunc(newEndpoint.path, newEndpoint.handleEndpoint).Methods(newEndpoint.method)

	if s.endpoints == nil {
		s.endpoints = make([]*endpoint, 0)
	}

	s.endpoints = append(s.endpoints, newEndpoint)

	return newEndpoint
}

func (s server) Verify() error {

	for _, endpoint := range s.endpoints {
		err := endpoint.check()
		if err != nil {
			return err
		}
	}

	return nil
}
