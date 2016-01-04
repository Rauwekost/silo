package command

import (
	log "github.com/rauwekost/silo/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/rauwekost/silo/Godeps/_workspace/src/github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	c := cobra.Command{
		Use: "silo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			v, _ := cmd.Flags().GetBool("verbose")
			if v == true {
				log.SetLevel(log.DebugLevel)
				log.Debug("Verbose logging enabled")
			}
		},
	}
	c.PersistentFlags().BoolP("verbose", "v", false, "More verbose output")
	c.AddCommand(NewVersionCommand())
	c.AddCommand(NewServerCommand())
	return &c
}
