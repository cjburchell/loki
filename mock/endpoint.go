package mock

import (
	"encoding/json"
	"fmt"
	"github.com/cjburchell/uatu-go"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

// endpoint configuration
type endpoint struct {
	// Request
	path   string
	method string

	// Stuff
	reply        reply
	handleReply  func(reply IReply, request *http.Request)
	handleVerify func(request *http.Request)
	name         string
	route        *mux.Route

	log log.ILog

	verify verify

	isVerbose bool
}

func (endpoint *endpoint) SetVerbose() IEndpoint {
	endpoint.isVerbose = true
	return endpoint
}

type verify struct {
	enabled       bool
	actualCalls   int
	expectedCalls int
}

func (verify *verify) handle(request *http.Request) {
	verify.actualCalls++
}

func createDefaultEndpoint(name, method, path string, log log.ILog) *endpoint {
	endpoint := &endpoint{name: name, path: path, method: method, log: log}
	endpoint.reply = reply{response: 200, log: log}

	verify := verify{}
	endpoint.verify = verify
	endpoint.handleVerify = verify.handle

	return endpoint
}

func (endpoint *endpoint) CustomVerify(handler func(request *http.Request)) IEndpoint {
	endpoint.handleVerify = handler
	return endpoint
}

func (endpoint *endpoint) Reply() IReply {
	return &endpoint.reply
}

func (endpoint *endpoint) CustomReply(handler func(reply IReply, request *http.Request)) IEndpoint {
	endpoint.handleReply = handler
	return endpoint
}

// IEndpoint interface
type IEndpoint interface {
	Reply() IReply
	CustomReply(handler func(reply IReply, request *http.Request)) IEndpoint
	CustomVerify(handler func(request *http.Request)) IEndpoint
	Once() IEndpoint
	Never() IEndpoint
	Times(count int) IEndpoint
	SetVerbose() IEndpoint
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
	endpoint.log.Printf("Handling endpoint %s %s %s", endpoint.name, endpoint.method, endpoint.path)

	var bodyString = ""
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err == nil {
		bodyString = string(bodyBytes)
	}

	requestData, _ := json.Marshal(struct {
		Endpoint    string `json:"endpoint"`
		Path        string `json:"path"`
		ContentType string `json:"content_type"`
		Body        string `json:"body"`
	}{
		Endpoint:    endpoint.name,
		Path:        endpoint.path,
		Body:        bodyString,
		ContentType: r.Header.Get("Content-type"),
	})

	if endpoint.isVerbose {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			endpoint.log.Error(err, "Unable to dump Request")
		}

		endpoint.log.Print(string(requestDump))

		vars := mux.Vars(r)
		if len(vars) != 0 {
			endpoint.log.Print("Values:")
			for key, value := range vars {
				endpoint.log.Printf("Key: %s, Value: $s", key, value)
			}
		}
	}

	currentReply := endpoint.reply
	if endpoint.handleReply != nil {
		currentReply = reply{response: 200}
		endpoint.handleReply(&currentReply, r)
	}

	currentReply.handle(w)
	fmt.Printf("Request:%s\n", string(requestData))

	endpoint.handleVerify(r)
}
