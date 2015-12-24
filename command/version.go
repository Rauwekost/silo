package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string
var Commit string

func NewVersionCommand() *cobra.Command {
	c := cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("v.%s %s\n", Version, Commit)
			return nil
		},
	}
	return &c
}
