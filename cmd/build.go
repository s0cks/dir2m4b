package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"charm.land/log/v2"
	"github.com/s0cks/dir2m4b/internal"
	"github.com/spf13/cobra"
)

func writeManifestToFile(data *internal.Metadata, targetPath string, outPath string) error {
	log.Info("writing manifest....", "target", targetPath, "out-path", outPath)

	manifest, err := internal.NewManifest(data, targetPath)
	if err != nil {
		return err
	}

	var manifestJson []byte
	manifestJson, err = json.MarshalIndent(manifest, "  ", "")
	if err != nil {
		return err
	}

	err = os.WriteFile(outPath, manifestJson, 0644)
	if err != nil {
		return err
	}

	log.Debug("done")
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build an m4b from the specified directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			builder := internal.NewBuilder()

			data, err := internal.NewMetadata(builder.SourceDir)
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

			outPath := getOutPath(builder.BuildDir)
			err = builder.Build(data, outPath)
			if err != nil {
				return err
			}

			err = writeManifestToFile(data, outPath, filepath.Join(builder.BuildDir, "manifest.json"))
			if err != nil {
				return err
			}

			return nil
		},
	}
	internal.AddBuilderFlags(cmd.Flags())
	cmd.Flags().String("spec", "", "The path of the m4b spec")

	cmd.Flags().String("title", "", "The title of the m4b")
	cmd.Flags().String("genre", "", "The genre of the m4b")
	cmd.Flags().String("artist", "", "The artist of the m4b")
	cmd.MarkFlagRequired("spec")
	RootCommand.AddCommand(cmd)
}
