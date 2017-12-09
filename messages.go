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
	Hostname               string    `json:"hostname"`
	PingTimeHumanFriendly  string    `json:"ping_time_humanfriendly"`
	LastAlertHumanFriendly string    `json:"last_alert_humanfriendly"`
	PingTime               time.Time `json:"ping_time"`
	LastAlert              time.Time `json:"last_alert"`
}

func notice(m pingMessage) error {
	h[m.Hostname] = &pingMessage{
		Hostname: m.Hostname,
		PingTime: time.Now(),
	}

	return nil
}
