package cmd

import (
	"path/filepath"

	"charm.land/log/v2"
	"github.com/s0cks/dir2m4b/internal"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build an m4b from the specified directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceDir, err := getSourceDir()
			if err != nil {
				return err
			}

			buildDir, err := getBuildDir(sourceDir)
			if err != nil {
				return err
			}

			log.Info("building....", "build-dir", buildDir, "source-dir", sourceDir)
			specPath, err := getSpecPath(sourceDir)
			if err != nil {
				return err
			}
			dataSpec, err := internal.LoadMetadataSpec(specPath)
			if err != nil {
				return err
			}

			data, err := internal.NewMetadata(sourceDir, &dataSpec)
			if err != nil {
				return err
			}

			dataPath := filepath.Join(buildDir, "Metadata")
			err = internal.WriteMetadataToFile(data, dataPath)
			if err != nil {
				return err
			}

			outPath := getOutPath(buildDir)
			err = data.Build(outPath, dataPath, "")
			if err != nil {
				return nil
			}

			return nil
		},
	}
	cmd.Flags().String("source-dir", "", "The source directory")
	cmd.Flags().String("build-dir", "", "The build directory")
	cmd.Flags().String("spec", "", "The path of the m4b spec")
	cmd.Flags().String("out", "", "The output filename")

	cmd.Flags().String("title", "", "The title of the m4b")
	cmd.Flags().String("genre", "", "The genre of the m4b")
	cmd.Flags().String("artist", "", "The artist of the m4b")
	cmd.MarkFlagRequired("spec")
	RootCommand.AddCommand(cmd)
}
