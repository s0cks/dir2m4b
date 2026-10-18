package cmd

import (
	"path/filepath"

	"charm.land/log/v2"
	"github.com/s0cks/dir2m4b/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cmd := &cobra.Command{
		Use:   "write-metadata",
		Short: "Write the m4b metadata file to the build directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceDir, err := getSourceDir()
			if err != nil {
				return err
			}

			buildDir, err := getBuildDir(sourceDir)
			if err != nil {
				return err
			}

			dataSpec := internal.MetadataSpec{
				Title:  viper.GetString("title"),
				Artist: viper.GetString("artist"),
				Genre:  viper.GetString("genre"),
			}

			log.Info("writing metadata....", "build-dir", buildDir, "source-dir", sourceDir)
			data, err := internal.NewMetadata(sourceDir, &dataSpec)
			if err != nil {
				return err
			}

			err = internal.WriteMetadataToFile(data, filepath.Join(buildDir, "Metadata"))
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().String("source-dir", "", "The source directory")
	cmd.Flags().String("build-dir", "", "The build directory")

	cmd.Flags().String("title", "", "The title of the m4b")
	cmd.Flags().String("genre", "", "The genre of the m4b")
	cmd.Flags().String("artist", "", "The artist of the m4b")
	RootCommand.AddCommand(cmd)
}
