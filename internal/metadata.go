package internal

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"charm.land/log/v2"
	"go.yaml.in/yaml/v3"
)

//go:embed Metadata
var MetadataTemplate string

type MetadataSpec struct {
	Title  string `json:"title" yaml:"title"`
	Cover  string `json:"cover" yaml:"cover"`
	Genre  string `json:"genre" yaml:"genre"`
	Artist string `json:"artist" yaml:"artist"`
}

type Metadata struct {
	MetadataSpec
	Root          string
	Chapters      []*Chapter
	TotalDuration int64 // milliseconds
}

func LoadMetadataSpec(source string) (MetadataSpec, error) {
	var spec MetadataSpec
	if source == "" {
		return spec, nil
	}

	if strings.HasSuffix(source, ".json") {
		bytes, err := os.ReadFile(source)
		if err != nil {
			return spec, fmt.Errorf("failed to read MetadataSpec from %s: %v", source, err)
		}

		if err := json.Unmarshal(bytes, &spec); err != nil {
			return spec, fmt.Errorf("failed to Unmarshal MetadataSpec from %s: %v", source, err)
		}
	} else if strings.HasSuffix(source, "jsonnet") {

	} else if strings.HasSuffix(source, ".yaml") || strings.HasSuffix(source, ".yml") {
		bytes, err := os.ReadFile(source)
		if err != nil {
			return spec, fmt.Errorf("failed to read MetadataSpec from %s: %v", source, err)
		}

		if err := yaml.Unmarshal(bytes, &spec); err != nil {
			return spec, fmt.Errorf("failed to Unmarhsal MetadataSpec from %s: %v", source, err)
		}
	} else {
		return MetadataSpec{}, fmt.Errorf("invalid metadata source: %s", source)
	}

	log.Info(fmt.Sprintf("loaded %v", spec))
	return spec, nil
}

func NewMetadata(dir string, spec *MetadataSpec) (*Metadata, error) {
	chapters, err := loadChaptersInDir(dir)
	if err != nil {
		return nil, err
	}

	var totalMs int64

	for _, ch := range chapters {
		ch.Start = totalMs
		totalMs = totalMs + ch.Duration
	}

	data := &Metadata{
		Root:          dir,
		Title:         spec.Title,
		Cover:         spec.Cover,
		Genre:         spec.Genre,
		Artist:        spec.Artist,
		Chapters:      chapters,
		TotalDuration: totalMs,
	}
	return data, nil
}

func WriteMetadataToFile(data *Metadata, dst string) error {
	log.Debug("writing metadata to file", "dest", dst, "title", data.Title, "grene")

	file, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	funcs := template.FuncMap{
		"add": func(a int64, b int64) int64 {
			return a + b
		},
	}

	tmpl, err := template.New("metadata").
		Funcs(funcs).
		Parse(MetadataTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	log.Debug("finished", "dest", dst, "title", data.Title)
	return nil
}

func (m *Metadata) Build(filename string, metadata string, cover string) error {
	args := []string{
		"-y", // Overwrite output file if it exists
	}

	numChapters := len(m.Chapters)
	var filterInputs []string

	for i, ch := range m.Chapters {
		args = append(args, "-i", ch.Source)
		filterInputs = append(filterInputs, fmt.Sprintf("[%d:a]", i))
	}

	currentIndex := numChapters
	metadataIndex := -1
	if metadata != "" {
		args = append(args, "-i", metadata)
		metadataIndex = currentIndex
		currentIndex++
	}

	coverIndex := -1
	if cover != "" {
		args = append(args, "-i", cover)
		coverIndex = currentIndex
		currentIndex++
	}

	concatFilter := fmt.Sprintf("%s concat=n=%d:v=0:a=1 [outa]", strings.Join(filterInputs, ""), numChapters)
	args = append(args, "-filter_complex", concatFilter)

	args = append(args, "-map", "[outa]")
	if cover != "" && coverIndex != -1 {
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

	args = append(args, filename)

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
