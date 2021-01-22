package editroute

import (
	"encoding/json"
	"github.com/cjburchell/loki/mock-server"
	"github.com/cjburchell/loki/models"
	log "github.com/cjburchell/uatu-go"
	"github.com/gorilla/mux"
	"net/http"
)

//Setup the edit mock route
func Setup(r *mux.Router, logger log.ILog, server mockserver.IServer) {
	route := r.PathPrefix("/@mock").Subrouter()

	route.HandleFunc("/endpoint", func(writer http.ResponseWriter, _ *http.Request) {
		handleGetEndpoints(writer, logger, server)
	}).Methods("GET")

	route.HandleFunc("/endpoint/{id}", func(writer http.ResponseWriter, request *http.Request) {
		handleGetEndpoint(writer, request, logger, server)
	}).Methods("GET")

	route.HandleFunc("/endpoint", func(writer http.ResponseWriter, request *http.Request) {
		handleCreateEndpoint(writer, request, logger, server)
	}).Methods("POST")

	route.HandleFunc("/endpoint/{id}", func(writer http.ResponseWriter, request *http.Request) {
		handleUpdateEndpoint(writer, request, logger, server)
	}).Methods("PUT")

	route.HandleFunc("/endpoint/{id}", func(writer http.ResponseWriter, request *http.Request) {
		handleDeleteEndpoint(writer, request, logger, server)
	}).Methods("DELETE")

	route.HandleFunc("/settings", func(writer http.ResponseWriter, _ *http.Request) {
		handleGetSettings(writer, logger, server)
	}).Methods("GET")

	route.HandleFunc("/settings", func(writer http.ResponseWriter, request *http.Request) {
		handleUpdateSettings(writer, request, logger, server)
	}).Methods("PUT")
}

func handleUpdateSettings(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockserver.IServer) {
	item := models.Settings{}
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		logger.Errorf(err, "Unmarshal Failed %s", request.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server.UpdateSettings(item)
	w.WriteHeader(http.StatusNoContent)
}

func handleGetSettings(w http.ResponseWriter, logger log.ILog, server mockserver.IServer) {
	settings := server.GetSettings()

	reply, _ := json.Marshal(settings)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write(reply)
	if err != nil {
		logger.Error(err)
	}
}

func handleGetEndpoints(w http.ResponseWriter, logger log.ILog, server mockserver.IServer) {
	endpoints := server.GetEndpoints()

	reply, _ := json.Marshal(endpoints)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write(reply)
	if err != nil {
		logger.Error(err)
	}
}

func handleGetEndpoint(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockserver.IServer) {
	vars := mux.Vars(request)
	endpointID := vars["id"]
	endpoint, err := server.GetEndpoint(endpointID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Error(err)
	}

	reply, _ := json.Marshal(endpoint)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(reply)
	if err != nil {
		logger.Error(err)
	}
}

func handleDeleteEndpoint(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockserver.IServer) {
	vars := mux.Vars(request)
	endpointID := vars["id"]
	err := server.DeleteEndpoint(endpointID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleCreateEndpoint(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockserver.IServer) {
	item := models.Endpoint{}
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		logger.Errorf(err, "Unmarshal Failed %s", request.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = server.AddEndpoint(item)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)

}

func handleUpdateEndpoint(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockserver.IServer) {
	vars := mux.Vars(request)
	endpointID := vars["id"]

	item := models.Endpoint{}
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		logger.Errorf(err, "Unmarshal Failed %s", request.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = server.UpdateEndpoint(endpointID, item)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
