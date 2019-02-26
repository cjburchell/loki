package mock

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cjburchell/go-uatu"

	"github.com/gorilla/mux"
)

// endpoint configuration
type endpoint struct {
	// Request
	path   string
	method string

	// Stuff
	reply        reply
	handleReply  func(writer http.ResponseWriter, request *http.Request)
	handleVerify func(writer http.ResponseWriter, request *http.Request)
	name         string
	route        *mux.Route

	verify verify
}

type verify struct {
	enabled       bool
	actualCalls   int
	expectedCalls int
}

func (verify *verify) handle(writer http.ResponseWriter, request *http.Request) {
	verify.actualCalls++
}

func createDefaultEndpoint(name, path, method string) *endpoint {
	endpoint := &endpoint{name: name, path: path, method: method}
	endpoint.reply = reply{response: 200}
	endpoint.handleReply = endpoint.reply.handle

	verify := verify{}
	endpoint.verify = verify
	endpoint.handleVerify = verify.handle

	return endpoint
}

func (endpoint *endpoint) CustomVerify(handler func(writer http.ResponseWriter, request *http.Request)) IEndpoint {
	endpoint.handleVerify = handler
	return endpoint
}

func (endpoint *endpoint) Reply() IReply {
	return &endpoint.reply
}

func (endpoint *endpoint) CustomReply(handler func(writer http.ResponseWriter, request *http.Request)) IEndpoint {
	endpoint.handleReply = handler
	return endpoint
}

type IEndpoint interface {
	Reply() IReply
	CustomReply(handler func(writer http.ResponseWriter, request *http.Request)) IEndpoint
	CustomVerify(handler func(writer http.ResponseWriter, request *http.Request)) IEndpoint
	Once() IEndpoint
	Never() IEndpoint
	Times(count int) IEndpoint
}

func (verify verify) check(endpoint endpoint, t *testing.T) {
	if !verify.enabled {
		assert.Equal(t, verify.expectedCalls, verify.actualCalls, fmt.Sprintf("Expected %d calls got %d calls in endpoint %s %s %s", verify.expectedCalls, verify.actualCalls, endpoint.name, endpoint.method, endpoint.path))
	}
}

func (endpoint endpoint) check(t *testing.T) {
	endpoint.verify.check(endpoint, t)
}

func (endpoint *endpoint) Once() IEndpoint {
	endpoint.verify.enabled = true
	endpoint.verify.expectedCalls = 1
	return endpoint
}

func (endpoint *endpoint) Never() IEndpoint {
	endpoint.verify.enabled = true
	endpoint.verify.expectedCalls = 0
	return endpoint
}

func (endpoint *endpoint) Times(count int) IEndpoint {
	endpoint.verify.enabled = true
	endpoint.verify.expectedCalls = count
	return endpoint
}

func (endpoint *endpoint) handleEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling endpoint %s %s %s", endpoint.name, endpoint.method, endpoint.path)

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

	endpoint.handleReply(w, r)
	endpoint.handleVerify(w, r)
}
