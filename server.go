package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func startServer() error {
	hostPort := fmt.Sprintf("%s:%d", host, port)
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/ping", pingGETHandler).Methods("GET")
	api.HandleFunc("/ping", pingPOSTHandler).Methods("POST")

	// every now and then, check status of nodes
	go checkUptime()

	// run http server
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

	if err = json.Unmarshal(bytes, &incoming); err != nil {
		logrus.Error(err)
		return
	}

	if err = notice(incoming); err != nil {
		logrus.Error(err)
	}

	w.WriteHeader(http.StatusOK)
	logrus.Debugf("incoming ping from %q", r.RemoteAddr)

	fmt.Fprint(w, "all ok\n")
}

func checkUptime() {
	for {
		// wait
		time.Sleep(2 * time.Second)

		// check hosts
		if len(knownHosts) == 0 {
			logrus.Error("no hosts found")
			continue
		}

		for key, pingMessage := range knownHosts {
			logrus.Debugf("node: %s, previous ping: %s", key, humanize.Time(pingMessage.PingTime))
		}

	}
}
