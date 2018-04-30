package main

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

const (
	alertFrequency = 5.00 // in minutes
	downtrigger    = 5.00 // in minutes
)

var mailserver smtpServer
var receiver string

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type smtpServer struct {
	host     string
	port     int
	password string
}

func (s *smtpServer) ServerName() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func setupAlerting(user string, pass string, server string) {
	if len(user) == 0 || len(pass) == 0 || len(server) == 0 {
		logrus.Error("email user and/or password empty! (or smtp server)")
	} else {
		logrus.Infof("enabling alerts for email %q", user)
		receiver = user
		mailserver = smtpServer{
			host:     server,
			password: pass,
			port:     465,
		}
	}
}

func notifyDowntime(m *pingMessage) error {

	mail := Mail{
		senderId: receiver,
		toIds:    []string{receiver},
		subject:  fmt.Sprintf("host %s is down", m.Hostname),
		body:     fmt.Sprintf("host down: %s\nlast seen: %s, \n\n bye!", m.Hostname, humanize.Time(m.pingTime)),
	}
	messageBody := mail.BuildMessage()

	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, mailserver.password, mailserver.host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         mailserver.host,
	}

	conn, err := tls.Dial("tcp", mailserver.ServerName(), tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, mailserver.host)
	if err != nil {
		return err
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		return err
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		return err
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			return err
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()
	logrus.Errorf("alert sent out for host %s", m.Hostname)

	return nil
}
