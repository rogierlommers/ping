package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	router.HandleFunc("/", rootHandler).Methods("GET")

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

	var response []pingMessage
	for _, message := range h {
		pm := pingMessage{
			Hostname:               message.Hostname,
			PingTime:               message.PingTime,
			PingTimeHumanFriendly:  humanize.Time(message.PingTime),
			LastAlertHumanFriendly: humanize.Time(message.LastAlert),
			LastAlert:              message.LastAlert,
		}

		response = append(response, pm)
	}

	b, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		logrus.Error(err)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(b)
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
		if len(h) == 0 {
			logrus.Error("no hosts found")
			continue
		}

		for key, m := range h {
			// calculate difference
			downtimeDuration := time.Since(m.PingTime)
			if downtimeDuration.Minutes() > downtrigger {
				// node considered down
				lastMailDuration := time.Since(m.LastAlert)
				if lastMailDuration.Minutes() > alertFrequency {
					notifyDowntime(m)
					m.LastAlert = time.Now()
				} else {
					secondsUntilAlert := alertFrequency - lastMailDuration.Minutes()
					logrus.Debugf("host %s down, last ping: %s, alerting in %f minutes", key, humanize.Time(m.PingTime), secondsUntilAlert)
				}
			}
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	io.WriteString(w, "<html><head><title>pingback</title></head><body><h2>pingback server</h2><div><a href='/api/ping'>/api/ping</a></div></body></html>")
}
