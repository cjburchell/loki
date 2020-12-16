package status

import (
	"github.com/cjburchell/uatu-go"
	"github.com/gorilla/mux"
	"net/http"
)

// Setup the status route
func Setup(r *mux.Router, logger log.ILog) {
	r.HandleFunc("/@status", func(writer http.ResponseWriter, request *http.Request) {
		handleGetStatus(writer, request, logger)
	}).Methods("GET")
}

func handleGetStatus(w http.ResponseWriter, r *http.Request, logger log.ILog) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("Ok"))
	if err != nil {
		logger.Error(err)
	}
}