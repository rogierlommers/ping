package main

import (
	"flag"
	"strings"

	"github.com/sirupsen/logrus"
)

// package globals
var (
	targetServer string
	mode         string
	host         string
	port         int
)

func init() {
	flag.StringVar(&mode, "mode", "client", "specify if we need to run as a server or as a client")
	flag.StringVar(&targetServer, "server", "http://localhost:8080", "the location of the server, f.e. http://ping.lommers.org")
	flag.StringVar(&host, "host", "localhost", "host to bind on, f.e. localhost")
	flag.IntVar(&port, "port", 8080, "port number to bind on, f.e. 8080")
	flag.Parse()

	// set log level
	logrus.SetLevel(logrus.DebugLevel)
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
