package mock

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"context"

	"github.com/cjburchell/go-uatu"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type IServer interface {
	Start(port int)
	Stop(t *testing.T)
	Endpoint(name, method, path string) IEndpoint
	Verify(t *testing.T)
}

type server struct {
	name      string
	r         *mux.Router
	srv       *http.Server
	endpoints []*endpoint
}

func CreateServer(name string) IServer {
	return &server{name: name}
}

func (s *server) Start(port int) {
	if s.r == nil {
		s.r = mux.NewRouter()
		log.Warnf("%s: Starting server with no endpoints", s.name)
	}

	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, s.r, func(writer io.Writer, params handlers.LogFormatterParams) {
		log.Printf("%s: \"%s %s\" Code:%d",
			s.name,
			params.Request.Method,
			params.URL.Path,
			params.StatusCode,
		)
	})

	s.srv = &http.Server{
		Handler:      loggedRouter,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("%s: Started http  server on port: %d", s.name, port)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()
}

func (s *server) Stop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error(err)
	}

	if t != nil {
		s.Verify(t)
	}
}

func (s *server) Endpoint(name, method, path string) IEndpoint {
	newEndpoint := createDefaultEndpoint(name, method, path)

	if s.r == nil {
		s.r = mux.NewRouter()
	}

	if len(newEndpoint.path) == 0 || len(newEndpoint.method) == 0 {
		return nil
	}

	log.Printf("%s: Loading newEndpoint %s %s %s", s.name, newEndpoint.name, newEndpoint.method, newEndpoint.path)
	newEndpoint.route = s.r.HandleFunc(newEndpoint.path, newEndpoint.handleEndpoint).Methods(newEndpoint.method)

	if s.endpoints == nil {
		s.endpoints = make([]*endpoint, 0)
	}

	s.endpoints = append(s.endpoints, newEndpoint)

	return newEndpoint
}

func (s server) Verify(t *testing.T) {

	for _, endpoint := range s.endpoints {
		endpoint.check(t)
	}
}
