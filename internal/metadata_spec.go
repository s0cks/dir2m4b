package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"charm.land/log/v2"
	"go.yaml.in/yaml/v3"
)

type MetadataSpec struct {
	Title  string `json:"title,omitempty" yaml:"title,omitempty"`
	Cover  string `json:"cover,omitempty" yaml:"cover,omitempty"`
	Genre  string `json:"genre,omitempty" yaml:"genre,omitempty"`
	Artist string `json:"artist,omitempty" yaml:"artist,omitempty"`
}

func LoadMetadataSpecFromYaml(source string) (MetadataSpec, error) {
	log.Info("loading MetadataSpec from yaml....", "source", source)
	if source == "" {
		return MetadataSpec{}, nil
	}

	bytes, err := os.ReadFile(source)
	if err != nil {
		return MetadataSpec{}, fmt.Errorf("failed to read MetadataSpec from %s: %v", source, err)
	}

	var spec MetadataSpec
	if err := yaml.Unmarshal(bytes, &spec); err != nil {
		return MetadataSpec{}, fmt.Errorf("failed to Unmarhsal MetadataSpec from %s: %v", source, err)
	}

	return spec, nil
}

func LoadMetadataSpecFromJson(source string) (MetadataSpec, error) {
	log.Info("loading MetadataSpec from json....", "source", source)
	bytes, err := os.ReadFile(source)
	if err != nil {
		return MetadataSpec{}, fmt.Errorf("failed to read MetadataSpec from %s: %v", source, err)
	}

	var spec MetadataSpec
	if err := json.Unmarshal(bytes, &spec); err != nil {
		return MetadataSpec{}, fmt.Errorf("failed to Unmarshal MetadataSpec from %s: %v", source, err)
	}

	return spec, nil
}

func LoadMetadataSpecFromFile(source string) (MetadataSpec, error) {
	if source == "" {
		return MetadataSpec{}, nil
	}

	if strings.HasSuffix(source, ".json") {
		return LoadMetadataSpecFromJson(source)
	} else if strings.HasSuffix(source, "jsonnet") {

	} else if strings.HasSuffix(source, ".yaml") || strings.HasSuffix(source, ".yml") {
		return LoadMetadataSpecFromYaml(source)
	}

	return MetadataSpec{}, fmt.Errorf("invalid metadata source: %s", source)
}
