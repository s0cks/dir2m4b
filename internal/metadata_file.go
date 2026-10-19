package internal

import (
	_ "embed"
	"os"
	"text/template"

	"charm.land/log/v2"
)

//go:embed Metadata
var MetadataTemplate string

func getDefaultTemplate() (*template.Template, error) {
	funcs := template.FuncMap{
		"add": func(a int64, b int64) int64 {
			return a + b
		},
	}
	tmpl, err := template.New("metadata").
		Funcs(funcs).
		Parse(MetadataTemplate)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

type MetadataWriter struct {
	TargetPath string
	Template   *template.Template
}

func NewMetadataWriter(target string) (*MetadataWriter, error) {
	tmpl, err := getDefaultTemplate()
	if err != nil {
		return nil, err
	}

	writer := MetadataWriter{
		TargetPath: target,
		Template:   tmpl,
	}
	return &writer, nil
}

func (writer *MetadataWriter) Write(data *Metadata) error {
	log.Info("writing metadata to file", "target", writer.TargetPath, "title", data.Title, "grene")
	file, err := os.OpenFile(writer.TargetPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := writer.Template.Execute(file, data); err != nil {
		return err
	}

	log.Info("finished", "target", writer.TargetPath, "title", data.Title)
	return nil
}
