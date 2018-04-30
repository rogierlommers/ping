package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
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

	srv := &http.Server{
		Handler:      router,
		Addr:         hostPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// run http server
	logrus.Infof("running as a server on %s", hostPort)
	return srv.ListenAndServe()
}

func pingGETHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	var response []pingMessage
	for _, message := range h {

		pm := pingMessage{
			Hostname:              message.Hostname,
			PingTimeHumanFriendly: humanize.Time(message.pingTime),
			IPv4:    message.IPv4,
			NoAlert: message.NoAlert,
		}

		if message.lastAlert.IsZero() {
			pm.LastAlertHumanFriendly = "never"
		} else {
			pm.LastAlertHumanFriendly = humanize.Time(message.lastAlert)
		}

		response = append(response, pm)
	}

	// now sort by hostname
	sort.Sort(messageSorter(response))

	// and send to client
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
			logrus.Error("did not receive a ping at all")
			continue
		}

		for key, m := range h {
			// calculate difference
			downtimeDuration := time.Since(m.pingTime)
			if downtimeDuration.Minutes() > downtrigger {
				// node considered down
				lastMailDuration := time.Since(m.lastAlert)
				if lastMailDuration.Minutes() > alertFrequency {
					if !m.NoAlert {
						logrus.Infof("triggering alert %v", m)
						notifyDowntime(m)
						m.lastAlert = time.Now()
					} else {
						logrus.Debugf("skip alert for host %s", m.Hostname)
					}
				} else {
					secondsUntilAlert := alertFrequency - lastMailDuration.Minutes()
					logrus.Infof("host %s down, last ping: %s, alerting in %f minutes", key, humanize.Time(m.pingTime), secondsUntilAlert)
				}
			}
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	io.WriteString(w, "<html><head><title>pingback</title></head><body><h2>pingback server</h2><div><a href='/api/ping'>/api/ping</a></div></body></html>")
}
