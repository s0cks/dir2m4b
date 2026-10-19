package internal

import (
	"github.com/spf13/viper"
	"path/filepath"
)

type Metadata struct {
	MetadataSpec
	Root          string     `json:"-" yaml:"-"`
	Chapters      []*Chapter `json:"-" yaml:"-"`
	TotalDuration int64      `json:"-" yaml:"-"`
	Source        string     `json:"-" yaml:"-"`
	Dest          string     `json:"-" yaml:"-"`
}

func findMetadataSpecFile(dir string) string {
	files := []string{
		"metadata.yml",
		"metadata.yaml",
		"metadata.json",
		"metadata.jsonnet",
		"spec.yml",
		"spec.yaml",
		"spec.json",
		"spec.jsonnet",
	}

	for _, file := range files {
		if fileExists(filepath.Join(dir, file)) {
			return filepath.Join(dir, file)
		}
	}

	return ""
}

func setCliOverride(name string, defaultValue string, out *string) {
	value := viper.GetString(name)
	if value != defaultValue {
		(*out) = value
	}
}

func NewMetadata(dir string) (*Metadata, error) {
	var spec MetadataSpec = MetadataSpec{}
	filename := findMetadataSpecFile(dir)
	if filename != "" {
		var err error
		spec, err = LoadMetadataSpecFromFile(filename)
		if err != nil {
			return nil, err
		}
	}

	setCliOverride("title", "", &spec.Title)
	setCliOverride("artist", "", &spec.Artist)
	setCliOverride("genre", "", &spec.Genre)
	setCliOverride("cover", "", &spec.Cover)

	chapters, err := loadChaptersInDir(dir)
	if err != nil {
		return nil, err
	}

	var totalMs int64

	for _, ch := range chapters {
		ch.Start = totalMs
		totalMs = totalMs + ch.Duration
	}

	data := Metadata{
		MetadataSpec:  spec,
		Root:          dir,
		Chapters:      chapters,
		TotalDuration: totalMs,
	}
	return &data, nil
}
