package main

import (
	"time"
)

var h knownHosts

type knownHosts map[string]*pingMessage

func init() {
	h = make(knownHosts)
}

type pingMessage struct {
	Hostname               string `json:"hostname"`
	PingTimeHumanFriendly  string `json:"last_ping"`
	LastAlertHumanFriendly string `json:"last_alert"`
	pingTime               time.Time
	lastAlert              time.Time
}

func notice(m pingMessage) error {
	if _, ok := h[m.Hostname]; ok {
		// known hostname, only update pingtime
		h[m.Hostname].pingTime = time.Now()
	} else {
		// unknown (new) hostname
		h[m.Hostname] = &pingMessage{
			Hostname: m.Hostname,
			pingTime: time.Now(),
		}
	}

	return nil
}
