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
	PingTime               time.Time
	LastAlert              time.Time
}

func notice(m pingMessage) error {
	if _, ok := h[m.Hostname]; ok {
		// known hostname, only update pingtime
		h[m.Hostname].PingTime = time.Now()
	} else {
		// unknown (new) hostname
		h[m.Hostname] = &pingMessage{
			Hostname: m.Hostname,
			PingTime: time.Now(),
		}
	}

	return nil
}
