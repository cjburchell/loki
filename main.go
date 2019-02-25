package main

import (
	"github.com/cjburchell/restmock/mock"
	"github.com/cjburchell/tools-go/env"

	"github.com/cjburchell/go-uatu"

	"github.com/cjburchell/restmock/config"
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

	server := mock.CreateServer()

	for _, endpoint := range endpoints {
		server.Endpoint(endpoint.Name, endpoint.Method, endpoint.Path).Reply().Body(endpoint.ResponseBody).Content(endpoint.ContentType).Code(endpoint.Response).FullHeader(endpoint.Header)
	}

	server.Start(port)

	// wait for ever
	ch := make(chan int)
	<-ch
}
