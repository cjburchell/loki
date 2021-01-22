package main

import (
	"context"
	"fmt"
	"github.com/cjburchell/loki/mock-server"
	"github.com/cjburchell/loki/models"
	"github.com/cjburchell/loki/routes/client-route"
	"github.com/cjburchell/loki/routes/edit-route"
	"github.com/cjburchell/loki/routes/status-route"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"time"

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

	startServer(port, endpoints, logger, set)
}

func startServer(port int, endpoints []models.Endpoint, log log.ILog, settings settings.ISettings) {
	r := mux.NewRouter()

	statusroute.Setup(r, log)

	serviceName := settings.Get("ServerName", "Loki")
	clientroute.Setup(r, settings.Get("ClientLocation", "client/dist/client"), log)

	ms := mockserver.Setup(serviceName,
		settings.GetInt("DefaultReply", http.StatusBadRequest),
		settings.Get("PartialMockServerAddress", ""), log)

	for _, endpointConfig := range endpoints {
		err := ms.AddEndpoint(endpointConfig)
		if err != nil {
			log.Errorf(err, "Unable to add endpoint %s", endpointConfig.Name)
		}
	}

	editroute.Setup(r, log, ms)
	ms.SetRoute(r)

	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, r, func(writer io.Writer, params handlers.LogFormatterParams) {
		log.Printf("%s: \"%s %s\" Code:%d",
			serviceName,
			params.Request.Method,
			params.URL.Path,
			params.StatusCode,
		)
	})

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("%s: Started http  server on port: %d", serviceName, port)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error(err)
		}
	}()

	// wait for ever
	ch := make(chan int)
	<-ch
}
