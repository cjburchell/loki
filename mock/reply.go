package mock

import (
	"encoding/json"
	"encoding/xml"
	log "github.com/cjburchell/uatu-go"
	"net/http"
	"time"
)

// IReply interface
type IReply interface {
	JSONBody(body interface{}) IReply
	Content(content string) IReply
	Code(code int) IReply
	Header(key, value string) IReply
	RawBody(body []byte) IReply
	StringBody(body string) IReply
	XMLBody(body interface{}) IReply
	FullHeader(header map[string]string) IReply
	Delay(delayTime int) IReply
}

type reply struct {
	responseBody []byte
	contentType  string
	response     int
	header       map[string]string
	delay        int
	log          log.ILog
}

func (reply reply) handle(w http.ResponseWriter) {

	if reply.delay != 0 {
		reply.log.Printf("Waiting for %sms", reply.delay)
		time.Sleep(time.Duration(reply.delay) * time.Millisecond)
	}

	reply.log.Printf("Send Response: %d %s Body: %s", reply.response, reply.contentType, reply.responseBody)

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
			reply.log.Error(err, "Unable to write response")
		}
	}
}

func (reply *reply) JSONBody(body interface{}) IReply {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		reply.log.Error(err, "Marshal Json")
	}
	return reply.RawBody(bodyBytes)
}

func (reply *reply) XMLBody(body interface{}) IReply {
	bodyBytes, err := xml.MarshalIndent(body, "  ", "    ")
	if err != nil {
		reply.log.Error(err, "Marshal Json")
	}

	return reply.RawBody([]byte(xml.Header + string(bodyBytes)))
}

func (reply *reply) RawBody(body []byte) IReply {
	reply.responseBody = body
	reply.log.Printf("Setting Reply Body of %s", reply.responseBody)
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

func (reply *reply) Delay(delayTime int) IReply {
	reply.delay = delayTime
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
