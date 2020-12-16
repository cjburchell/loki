package clientroute

import (
	"mime"
	"net/http"

	"github.com/cjburchell/uatu-go"
	"github.com/gorilla/mux"
)

func Setup(r *mux.Router, clientLocation string, logger log.ILog) {
	route :=  r.PathPrefix("/@client").Subrouter()

	err := mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
	if err != nil {
		logger.Error(err)
	}

	err = mime.AddExtensionType(".html", "text/html; charset=utf-8")
	if err != nil {
		logger.Error(err)
	}

	handleClient := func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, clientLocation+"/index.html")
	}

	route.HandleFunc("/", handleClient)
	route.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(clientLocation))))
}
