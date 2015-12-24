package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

func NewVersionCommand() *cobra.Command {
	c := cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("v.%s\n", Version)
			return nil
		},
	}
	return &c
}
