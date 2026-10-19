package internal

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/log/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func AddBuilderFlags(flags *pflag.FlagSet) {
	flags.String("build-dir", "", "The build directory")
	flags.String("source-dir", "", "The source directory")
	flags.String("target", "", "The output filename")
	flags.String("cover", "", "The cover filename")
	flags.String("metadata", "", "The Metadata filename")
}

type Builder struct {
	BuildDir     string
	SourceDir    string
	TargetPath   string
	CoverPath    string
	MetadataPath string
}

func NewBuilder() Builder {
	return Builder{
		BuildDir:     viper.GetString("build-dir"),
		SourceDir:    viper.GetString("source-dir"),
		TargetPath:   viper.GetString("target"),
		CoverPath:    viper.GetString("cover"),
		MetadataPath: viper.GetString("metadata"),
	}
}

func (builder *Builder) createFfmpegBuildArgs(mdata *Metadata) []string {
	args := []string{
		"-y", // Overwrite output file if it exists
	}

	numChapters := len(mdata.Chapters)
	var filterInputs []string

	for i, ch := range mdata.Chapters {
		args = append(args, "-i", ch.Source)
		filterInputs = append(filterInputs, fmt.Sprintf("[%d:a]", i))
	}

	currentIndex := numChapters
	metadataIndex := -1
	if builder.MetadataPath != "" {
		args = append(args, "-i", builder.MetadataPath)
		metadataIndex = currentIndex
		currentIndex++
	}

	coverIndex := -1
	if builder.CoverPath != "" {
		args = append(args, "-i", builder.CoverPath)
		coverIndex = currentIndex
		currentIndex++
	}

	concatFilter := fmt.Sprintf("%s concat=n=%d:v=0:a=1 [outa]", strings.Join(filterInputs, ""), numChapters)
	args = append(args, "-filter_complex", concatFilter)

	args = append(args, "-map", "[outa]")
	if builder.CoverPath != "" && coverIndex != -1 {
		args = append(args, "-map", fmt.Sprintf("%d:v", coverIndex))
		args = append(args, "-c:v", "mjpeg", "-disposition:v", "attached_pic")
	}

	if metadataIndex != -1 {
		args = append(args, "-map_metadata", fmt.Sprintf("%d", metadataIndex))
	}

	args = append(args,
		"-c:a", "aac",
		"-b:a", "64k",
	)

	args = append(args, builder.TargetPath)
	return args
}

func (builder *Builder) Build(mdata *Metadata, filename string) error {
	if builder.MetadataPath == "" {
		mdataFilename := filepath.Join(builder.BuildDir, "Metadata")
		if fileExists(mdataFilename) {
			builder.MetadataPath = mdataFilename
		}
	}

	if builder.TargetPath == "" {
		builder.TargetPath = filepath.Join(builder.BuildDir, "out.m4b")
	}

	args := builder.createFfmpegBuildArgs(mdata)
	log.Info(fmt.Sprintf("running ffmpeg %s", strings.Join(args, " ")))
	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute ffmpeg: %v, stderr: %s", err, stderr.String())
	}

	return nil
}
