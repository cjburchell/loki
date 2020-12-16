package mockServer

import (
	"fmt"
	"github.com/cjburchell/loki/models"
	"github.com/cjburchell/uatu-go"
	"github.com/gorilla/mux"
	"net/http"
)

// IServer interface
type IServer interface {
	AddEndpoint(endpoint models.Endpoint) error
	DeleteEndpoint(id string) error
	UpdateEndpoint(id string, endpoint models.Endpoint) error
	GetEndpoint(id string) (*models.Endpoint, error)
	GetEndpoints() []models.Endpoint
}

type server struct {
	name                     string
	defaultReply             int
	partialMockServerAddress string
	client                   *http.Client
	endpoints                []*endpoint
	log                      log.ILog
}

func (s *server) DeleteEndpoint(id string) error {
	for i, ep := range s.endpoints {
		if ep.Name == id{
			// delete endpoint from list
			s.endpoints = append(s.endpoints[:i], s.endpoints[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("unable to find endpoint %s", id)
}

func (s *server) UpdateEndpoint(id string, endpoint models.Endpoint) error {

	if len(endpoint.Path) == 0  {
		return fmt.Errorf("missing Path")
	}

	if len(endpoint.Method) == 0 {
		return fmt.Errorf("missing Method")
	}

	if len(endpoint.Name) == 0  {
		return fmt.Errorf("missing endpint name")
	}

	if endpoint.Name != id  {
		return fmt.Errorf("name mismach %s != %s", id , endpoint.Name)
	}

	ep, err := s.getEndpoint(id)
	if err != nil{
		return err
	}

	ep.Endpoint = endpoint
	ep.updateResponseBody()
	return nil
}

func (s server) getEndpoint(id string) (*endpoint, error) {
	for _, ep := range s.endpoints {
		if ep.Name == id{
			return ep, nil
		}
	}

	return nil, fmt.Errorf("unable to find endpoint %s", id)
}

func (s server) GetEndpoint(id string) (*models.Endpoint, error) {
	ep, err := s.getEndpoint(id)
	if err != nil{
		return nil, err
	}

	return &ep.Endpoint, nil
}

func (s server) GetEndpoints() []models.Endpoint {
	endpoints := make([]models.Endpoint, 0)
	for index, ep := range s.endpoints {
		endpoints[index] = ep.Endpoint
	}
	return endpoints
}

// CreateServer creates the server
func Setup(name string, defaultReply int, partialMockServerAddress string, log log.ILog, r *mux.Router) IServer {
	server := &server{
		name:                     name,
		defaultReply:             defaultReply,
		partialMockServerAddress: partialMockServerAddress,
		client:                   &http.Client{},
		log: log,
	}

	r.PathPrefix("/").HandlerFunc(server.defaultHandler)

	return server
}

func (s *server) defaultHandler(w http.ResponseWriter, r *http.Request) {

	// handle mocked endpoints
	for _, endpoint := range s.endpoints {
		if endpoint.Method == r.Method && endpoint.Path == r.URL.Path {
			endpoint.handleEndpoint(w, r)
			return
		}
	}

	if len(s.partialMockServerAddress) == 0 { // no partial mock
	    // just handle default stuff
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

func (s *server) AddEndpoint(ep models.Endpoint) error {

	if len(ep.Path) == 0  {
		return fmt.Errorf("missing Path")
	}

	if len(ep.Method) == 0 {
		return fmt.Errorf("missing Method")
	}

	if len(ep.Name) == 0  {
		return fmt.Errorf("missing endpint name")
	}

	newEndpoint := createDefaultEndpoint(ep, s.log)
	s.log.Printf("%s: Loading newEndpoint %s %s %s", s.name, newEndpoint.Name, newEndpoint.Method, newEndpoint.Path)

	if s.endpoints == nil {
		s.endpoints = make([]*endpoint, 0)
	}

	s.endpoints = append(s.endpoints, newEndpoint)

	return nil
}
