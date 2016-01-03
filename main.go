package main

import (
	log "github.com/rauwekost/silo/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/rauwekost/silo/command"
)

var version string

func main() {
	command.Version = version
	rootCmd := command.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
