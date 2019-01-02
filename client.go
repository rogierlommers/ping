package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var client = &http.Client{}

func startClient() error {
	logrus.Infof("running as a client, connecting to %s", targetServer)

	message, err := newPingMessage()
	if err != nil {
		logrus.Errorf("unable to build ping message, abort ping: %q", err.Error())
		return nil
	}

	if err := pingBack(message); err != nil {
		logrus.Error(err)
	}

	return nil
}

func pingBack(message pingMessage) error {
	url := fmt.Sprintf("%s/api/ping", targetServer)
	logrus.Debugf("ping back url: %q", url)

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	logrus.Debugf("ping done, server responded %s", resp.Status)
	return nil
}

func newPingMessage() (pingMessage, error) {
	message := pingMessage{
		Hostname: getHostname(),
		IPv4:     getExternalIP(),
		NoAlert:  noAlert,
	}

	return message, nil
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return hostname
}

func getExternalIP() string {
	targetURL := "https://myexternalip.com/raw"

	resp, err := http.Get(targetURL)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}

	return strings.TrimSuffix(string(body), "\n")
}
