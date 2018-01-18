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

// messageSorter sorts pingMessages by hostname
type messageSorter []pingMessage

func (a messageSorter) Len() int           { return len(a) }
func (a messageSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a messageSorter) Less(i, j int) bool { return a[i].Hostname < a[j].Hostname }

// notice processes an incoming message
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
