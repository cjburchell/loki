package mock

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/cjburchell/go-uatu"
)

type IReply interface {
	JsonBody(body interface{}) IReply
	Content(content string) IReply
	Code(code int) IReply
	Header(key, value string) IReply
	RawBody(body []byte) IReply
	StringBody(body string) IReply
	XmlBody(body interface{}) IReply
	FullHeader(header map[string]string) IReply
}

type reply struct {
	responseBody []byte
	contentType  string
	response     int
	header       map[string]string
}

func (reply reply) handle(w http.ResponseWriter) {
	log.Printf("Send Response: %d %s Body: %s", reply.response, reply.contentType, reply.responseBody)

	if reply.header != nil {
		for key, value := range reply.header {
			w.Header().Set(key, value)
		}
	}

	if len(reply.contentType) != 0 {
		w.Header().Set("Content-Type", reply.contentType)
	}

	w.WriteHeader(reply.response)

	if reply.responseBody != nil {
		_, err := w.Write(reply.responseBody)
		if err != nil {
			log.Error(err, "Unable to write response")
		}
	}
}

func (reply *reply) JsonBody(body interface{}) IReply {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Error(err, "Marshal Json")
	}
	return reply.RawBody(bodyBytes)
}

func (reply *reply) XmlBody(body interface{}) IReply {
	bodyBytes, err := xml.MarshalIndent(body, "  ", "    ")
	if err != nil {
		log.Error(err, "Marshal Json")
	}

	return reply.RawBody([]byte(xml.Header + string(bodyBytes)))
}

func (reply *reply) RawBody(body []byte) IReply {
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
