package main

import "github.com/davecgh/go-spew/spew"
import "time"

var knownHosts map[string]pingMessage

func init() {
	knownHosts = make(map[string]pingMessage)
}

type pingMessage struct {
	Hostname string `json:"hostname"`
	PingTime time.Time
}

func notice(m pingMessage) error {

	knownHosts[m.Hostname] = pingMessage{
		Hostname: m.Hostname,
		PingTime: time.Now(),
	}

	// check if h is known

	// if h is not known, then add as new host

	// if h is known, then update time
	spew.Dump(knownHosts)
	return nil
}
