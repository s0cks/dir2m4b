package cmd

import (
	"fmt"

	"github.com/s0cks/dir2m4b/internal"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(internal.GetVersion())
			return nil
		},
	}
	RootCommand.AddCommand(cmd)
}
