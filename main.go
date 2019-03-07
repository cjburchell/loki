package main

import (
	"encoding/json"

	"github.com/cjburchell/loki/mock"
	"github.com/cjburchell/tools-go/env"

	"github.com/cjburchell/go-uatu"

	"github.com/cjburchell/loki/config"
)

func main() {
	configFile := env.Get("CONFIG_FILE", "config.json")

	err := config.Setup(configFile)
	if err != nil {
		log.Fatal(err, "Unable to setup config File")
		return
	}

	port := env.GetInt("PORT", 8080)

	endpoints, err := config.GetEndpoints()
	if err != nil {
		log.Fatal(err, "Unable build endpoints")
		return
	}

	startHTTPTestEndpoints(port, endpoints)
}

func startHTTPTestEndpoints(port int, endpoints []config.Endpoint) {

	server := mock.CreateServer(env.Get("SERVER_NAME", "Loki"))

	for _, endpointConfig := range endpoints {
		endpoint := server.Endpoint(endpointConfig.Name, endpointConfig.Method, endpointConfig.Path)
		reply := endpoint.Reply()

		if endpointConfig.ContentType == "application/json" {
			reply.JsonBody(endpointConfig.ResponseBody)
		} else {
			var body string
			err := json.Unmarshal(endpointConfig.ResponseBody, &body)
			if err == nil {
				reply.StringBody(body)
			}
		}

		reply.Content(endpointConfig.ContentType).Code(endpointConfig.Response).FullHeader(endpointConfig.Header)
	}

	server.Start(port)
	defer server.Stop(nil)

	// wait for ever
	ch := make(chan int)
	<-ch
}
