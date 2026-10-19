package cmd

import (
	"path/filepath"

	"github.com/s0cks/dir2m4b/internal"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "write-metadata",
		Short: "Write the m4b metadata file to the build directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			builder := internal.NewBuilder()

			data, err := internal.NewMetadata(builder.BuildDir)
			if err != nil {
				return err
			}

			writer, err := internal.NewMetadataWriter(filepath.Join(builder.BuildDir, "Metadata"))
			if err != nil {
				return err
			}

			if err = writer.Write(data); err != nil {
				return err
			}

			return nil
		},
	}
	internal.AddBuilderFlags(cmd.Flags())
	cmd.Flags().String("title", "", "The title of the m4b")
	cmd.Flags().String("genre", "", "The genre of the m4b")
	cmd.Flags().String("artist", "", "The artist of the m4b")
	RootCommand.AddCommand(cmd)
}
