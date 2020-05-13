package mock

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"context"

	"github.com/cjburchell/uatu-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// IServer interface
type IServer interface {
	Start(port int)
	Stop(t *testing.T)
	Endpoint(name, method, path string) IEndpoint
	Verify(t *testing.T)
}

type server struct {
	name                     string
	defaultReply             int
	partialMockServerAddress string
	client                   *http.Client
	r                        *mux.Router
	srv                      *http.Server
	endpoints                []*endpoint
	log                      log.ILog
}

// CreateServer creates the server
func CreateServer(name string, defaultReply int, partialMockServerAddress string, log log.ILog) IServer {
	return &server{
		name:                     name,
		defaultReply:             defaultReply,
		partialMockServerAddress: partialMockServerAddress,
		client:                   &http.Client{},
		log: log,
	}
}

func (s *server) Start(port int) {
	if s.r == nil {
		s.r = mux.NewRouter()
		s.log.Warnf("%s: Starting server with no endpoints", s.name)
	}

	s.r.PathPrefix("/").HandlerFunc(s.DefaultHandler)

	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, s.r, func(writer io.Writer, params handlers.LogFormatterParams) {
		s.log.Printf("%s: \"%s %s\" Code:%d",
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

	s.log.Printf("%s: Started http  server on port: %d", s.name, port)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()
}

func (s *server) DefaultHandler(w http.ResponseWriter, r *http.Request) {

	if len(s.partialMockServerAddress) == 0 { // no partial mock
		w.WriteHeader(s.defaultReply)
		return
	}

	req, err := http.NewRequest(r.Method, s.partialMockServerAddress+r.URL.Path, r.Body)
	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	resp, err := s.client.Do(req)
	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil{
			s.log.Error(err)
		}
	}()

	for key, values := range resp.Header {
		for _, value := range values{
			w.Header().Set(key, value)
		}
	}

	var body []byte
	_, err = r.Body.Read(body)
	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *server) Stop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Error(err)
	}

	if t != nil {
		s.Verify(t)
	}
}

func (s *server) Endpoint(name, method, path string) IEndpoint {
	newEndpoint := createDefaultEndpoint(name, method, path, s.log)

	if s.r == nil {
		s.r = mux.NewRouter()
	}

	if len(newEndpoint.path) == 0 || len(newEndpoint.method) == 0 {
		return nil
	}

	s.log.Printf("%s: Loading newEndpoint %s %s %s", s.name, newEndpoint.name, newEndpoint.method, newEndpoint.path)
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
