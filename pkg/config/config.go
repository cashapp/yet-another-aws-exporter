package config

import (
	"errors"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigPath       string
	DisabledScrapers []string `yaml:"disabledScrapers"`
}

func (c *Config) Load() error {
	if _, err := os.Stat(c.ConfigPath); errors.Is(err, os.ErrNotExist) {
		log.Debug("Config file does not exist! Moving on")
		return nil
	}

	yamlFile, err := ioutil.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	log.Debug("Config file loaded!")
	return nil
}
