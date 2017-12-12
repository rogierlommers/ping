package main

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	gmail "github.com/rogierlommers/go-gmail-sender"
	"github.com/sirupsen/logrus"
)

const (
	alertFrequency = 5.00 // in minutes
	downtrigger    = 5.00 // in minutes
)

var (
	mail     gmail.Client
	receiver string
)

func setupAlerting(user string, pass string) {
	if len(user) == 0 || len(pass) == 0 {
		logrus.Error("email user and/or password empty!")
	} else {
		logrus.Infof("enabling alerts for email %q", user)
		mail = gmail.NewClient(user, pass)
		receiver = user
	}
}

func notifyDowntime(m *pingMessage) {
	emailMessage := gmail.Message{
		Receiver: receiver,
		Subject:  fmt.Sprintf("host %s is down", m.Hostname),
		Body:     fmt.Sprintf("host down: %s\nlast seen: %s, \n\n bye!", m.Hostname, humanize.Time(m.PingTime)),
	}

	// finally send message
	err := mail.Send(emailMessage)
	if err != nil {
		logrus.Errorf("error sending alert: %q", err.Error())
		return
	}

	logrus.Errorf("alert sent out for host %s", m.Hostname)
}
