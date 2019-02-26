package mock

import (
	"encoding/json"
	"net/http"

	"github.com/cjburchell/go-uatu"
)

type IReply interface {
	Body(body interface{}) IReply
	Content(content string) IReply
	Code(code int) IReply
	Header(key, value string) IReply
	RawBody(body json.RawMessage) IReply
	StringBody(body string) IReply
	FullHeader(header map[string]string) IReply
}

type reply struct {
	responseBody json.RawMessage
	contentType  string
	response     int
	header       map[string]string
}

func (reply *reply) handle(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Send Response: %d %s Body: %s", reply.response, reply.contentType, reply.responseBody)
	w.WriteHeader(reply.response)
	w.Header().Set("Content-Type", reply.contentType)
	if reply.header != nil {
		log.Print("Header")
		for key, value := range reply.header {
			log.Printf("%s %s", key, value)
			w.Header().Set(key, value)
		}
	}

	var err error
	if reply.contentType == "application/json" {
		_, err = w.Write(reply.responseBody)
	} else {
		var body string
		err := json.Unmarshal(reply.responseBody, &body)
		if err == nil {
			_, err = w.Write([]byte(body))
		}
	}

	if err != nil {
		log.Error(err, "Unable to write response")
	}
}

func (reply *reply) Body(body interface{}) IReply {
	bodyBytes, _ := json.Marshal(body)
	return reply.RawBody(bodyBytes)
}

func (reply *reply) RawBody(body json.RawMessage) IReply {
	reply.responseBody = body
	log.Printf("Setting Reply Body of %s", reply.responseBody)
	return reply
}

func (reply *reply) StringBody(body string) IReply {
	return reply.RawBody([]byte(body))
}

func (reply *reply) Content(content string) IReply {
	reply.contentType = content
	return reply
}

func (reply *reply) Code(code int) IReply {
	reply.response = code
	return reply
}

func (reply *reply) Header(key, value string) IReply {
	if reply.header == nil {
		reply.header = map[string]string{}
	}

	reply.header[key] = value
	return reply
}

func (reply *reply) FullHeader(header map[string]string) IReply {
	reply.header = header
	return reply
}
