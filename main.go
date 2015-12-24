package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rauwekost/silo/command"
)

var version = "1.0"
var commit = "abc"

func main() {
	command.Version = version
	command.Commit = commit
	rootCmd := command.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
