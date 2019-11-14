package main

import (
    log "github.com/sirupsen/logrus"
    "ssh-to-k8s/cmd"
)

func init() {
    log.SetFormatter(&log.TextFormatter{
    })
    log.SetLevel(log.DebugLevel)
}

func main() {
    cmd.Execute()
}

