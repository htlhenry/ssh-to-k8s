package main

import (
	log "github.com/sirupsen/logrus"
	"ssh-to-k8s/cmd"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	cmd.Execute()
}
