package models

import "encoding/json"

//Endpoint model
type Endpoint struct {
	Path         string            `json:"path"`
	Method       string            `json:"method"`
	ResponseBody json.RawMessage   `json:"response_body"`
	StringBody   string            `json:"string_body"`
	ContentType  string            `json:"content_type"`
	Response     int               `json:"response"`
	Header       map[string]string `json:"header"`
	Headers      []Header          `json:"headers"`
	Name         string            `json:"name"`
	ReplyDelay   int               `json:"reply_delay"`
}

//Header item for the endpoint model
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
