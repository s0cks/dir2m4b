package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func getSourceDir() (string, error) {
	sourceDir := viper.GetString("source-dir")
	if sourceDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		sourceDir = cwd
	}

	return sourceDir, nil
}

func getSpecPath(sourceDir string) (string, error) {
	spec := viper.GetString("spec")
	if spec == "" {
		return "", fmt.Errorf("invalid spec path: %s", spec)
	}

	return spec, nil
}

func getOutPath(buildDir string) string {
	outPath := viper.GetString("out")
	if outPath == "" {
		return filepath.Join(buildDir, "out.m4b")
	}

	if strings.HasPrefix(outPath, buildDir) {
		outPath = filepath.Join(buildDir, outPath)
	}

	return outPath
}

func getBuildDir(sourceDir string) (string, error) {
	buildDir := viper.GetString("build-dir")
	if buildDir == "" {
		buildDir = filepath.Join(sourceDir, "build")
	}

	err := os.MkdirAll(buildDir, 0755)
	if err != nil {
		return "", err
	}

	return buildDir, nil
}
