package dtos

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigData struct {
	Gormite struct {
		Orm struct {
			Mapping map[string]*ConfigDataMapping
		}
	}
}

type ConfigDataMapping struct {
	Dir string
}

func NewConfigData(path string) (*ConfigData, error) {
	configContent, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	config := ConfigData{}

	if err = yaml.Unmarshal(configContent, &config); err != nil {
		return nil, errors.WithStack(err)
	}

	return &config, nil
}
