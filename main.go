package main

import (
	"encoding/json"
	"net/http"

	"github.com/cjburchell/loki/mock"
	"github.com/cjburchell/settings-go"
	"github.com/cjburchell/tools-go/env"
	log "github.com/cjburchell/uatu-go"
	logSettings "github.com/cjburchell/uatu-go/settings"

	"github.com/cjburchell/loki/config"
)

func main() {
	set := settings.Get(env.Get("SettingsFile", "settings.yaml"))
	logger := log.Create(logSettings.Get(set.GetSection("Logging")))

	configFile := set.Get("ConfigFile", "config.json")

	err := config.Setup(configFile)
	if err != nil {
		logger.Fatal(err, "Unable to setup config File")
		return
	}

	port := set.GetInt("Port", 8080)

	endpoints, err := config.GetEndpoints(logger)
	if err != nil {
		logger.Fatal(err, "Unable build endpoints")
		return
	}

	startHTTPTestEndpoints(port, endpoints, logger, set)
}

func startHTTPTestEndpoints(port int, endpoints []config.Endpoint, log log.ILog, settings settings.ISettings) {

	server := mock.CreateServer(settings.Get("ServerName", "Loki"),
		settings.GetInt("DefaultReply", http.StatusBadRequest),
		settings.Get("PartialMockServerAddress", ""), log)

	for _, endpointConfig := range endpoints {
		endpoint := server.Endpoint(endpointConfig.Name, endpointConfig.Method, endpointConfig.Path)
		reply := endpoint.Reply()

		if endpointConfig.ResponseBody != nil {
			if endpointConfig.ContentType == "application/json" {
				reply.JSONBody(endpointConfig.ResponseBody)
			} else {
				var body string
				err := json.Unmarshal(endpointConfig.ResponseBody, &body)
				if err == nil {
					reply.StringBody(body)
				}
			}
		}

		if len(endpointConfig.StringBody) != 0 {
			reply.StringBody(endpointConfig.StringBody)
		}

		reply.Content(endpointConfig.ContentType).Code(endpointConfig.Response).FullHeader(endpointConfig.Header).Delay(endpointConfig.ReplyDelay)
	}

	server.Start(port)
	defer server.Stop(nil)

	// wait for ever
	ch := make(chan int)
	<-ch
}
