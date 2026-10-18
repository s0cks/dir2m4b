package cmd

import (
	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build an m4b from the specified directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("building....")
			return nil
		},
	}
	RootCommand.AddCommand(cmd)
}
