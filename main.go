package main

import (
	"flag"
	"strings"

	"github.com/sirupsen/logrus"
)

// package globals
var (
	targetServer  string
	mode          string
	host          string
	emailUser     string
	emailPassword string
	port          int
	debug         bool
)

func init() {
	flag.StringVar(&mode, "mode", "client", "specify if we need to run as a server or as a client")
	flag.StringVar(&targetServer, "server", "http://localhost:8080", "the location of the server, f.e. http://ping.lommers.org")
	flag.StringVar(&host, "host", "localhost", "host to bind on, f.e. localhost")
	flag.StringVar(&emailUser, "emailuser", "", "gmail address to use for alerting")
	flag.StringVar(&emailPassword, "emailpassword", "", "gmail password to use for alerting")
	flag.BoolVar(&debug, "debug", false, "true for debug mode")
	flag.IntVar(&port, "port", 8080, "port number to bind on, f.e. 8080")
	flag.Parse()

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// setup alerting
	setupAlerting(emailUser, emailPassword)
}

func main() {
	switch strings.ToLower(mode) {

	case "client":
		if err := startClient(); err != nil {
			logrus.Error(err)
		}

	case "server":
		if err := startServer(); err != nil {
			logrus.Error(err)
		}

	default:
		logrus.Error("invalid mode provided")
	}

}
