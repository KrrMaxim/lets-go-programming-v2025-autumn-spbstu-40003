package config

import (
	"fmt"
	"os"
	"errors"
	"gopkg.in/yaml.v3"
)

var(
	ErrOpenFile   = errors.New("error: problem with opening file")
	ErrDecodeFile = errors.New("error: failed to decode file")
	ErrFileNotSet = errors.New("error: both input-file and output-file must be set")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ConfigLoad(path string) (*Config, error) {
	var config Config

	data, err := os.Open(path)
	if err != nil {
		return nil, ErrOpenFile
	}

	defer func(data *os.File) {
		err := data.Close()
		if err != nil {
			panic(fmt.Errorf("error: problem with closing file: %w", err))
		}
	}(data)

	decoder := yaml.NewDecoder(data)

	if err := decoder.Decode(&config); err != nil {
		return nil, ErrDecodeFile
	}

	if config.InputFile == "" || config.OutputFile == "" {
		return nil, ErrFileNotSet
	}

	return &config, nil
}
