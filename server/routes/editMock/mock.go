package editMock

import (
	"encoding/json"
	"github.com/cjburchell/loki/models"
	"github.com/cjburchell/loki/routes/mockServer"
	log "github.com/cjburchell/uatu-go"
	"github.com/gorilla/mux"
	"net/http"
)

func Setup(r *mux.Router, logger log.ILog, server mockServer.IServer) {
	route :=  r.PathPrefix("/@mock").Subrouter()

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
}

func handleGetEndpoints(w http.ResponseWriter, logger log.ILog, server mockServer.IServer) {
	endpoints := server.GetEndpoints()

	reply, _ := json.Marshal(endpoints)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write(reply)
	if err != nil {
		logger.Error(err)
	}
}

func handleGetEndpoint(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockServer.IServer) {
	vars := mux.Vars(request)
	endpointId := vars["id"]
	endpoint, err := server.GetEndpoint(endpointId)
	if err != nil{
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

func handleDeleteEndpoint(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockServer.IServer) {
	vars := mux.Vars(request)
	endpointId := vars["id"]
	err := server.DeleteEndpoint(endpointId)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		logger.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleCreateEndpoint(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockServer.IServer) {
	item := models.Endpoint{}
	err := json.NewDecoder(request.Body).Decode(item)
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

func handleUpdateEndpoint(w http.ResponseWriter, request *http.Request, logger log.ILog, server mockServer.IServer) {
	vars := mux.Vars(request)
	endpointId := vars["id"]

	item := models.Endpoint{}
	err := json.NewDecoder(request.Body).Decode(item)
	if err != nil {
		logger.Errorf(err, "Unmarshal Failed %s", request.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = server.UpdateEndpoint(endpointId, item)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}


