package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCommand = &cobra.Command{
	Use: "dir2m4b",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "-"))
		viper.SetEnvPrefix("DIR2M4B")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
