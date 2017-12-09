package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func startServer() error {
	hostPort := fmt.Sprintf("%s:%d", host, port)
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/ping", pingGETHandler).Methods("GET")
	api.HandleFunc("/ping", pingPOSTHandler).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         hostPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("running as a server on %s", hostPort)
	return srv.ListenAndServe()
}

func pingGETHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "all ok\n")
}

func pingPOSTHandler(w http.ResponseWriter, r *http.Request) {
	var incoming pingMessage

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err)
	}

	err = json.Unmarshal(bytes, &incoming)
	spew.Dump(incoming)

	w.WriteHeader(http.StatusOK)
	logrus.Debugf("incoming ping from %q", r.RemoteAddr)

	fmt.Fprint(w, "all ok\n")
}
