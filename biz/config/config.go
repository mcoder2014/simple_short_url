package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL      string `yaml:"base_url"`
	ShortURLFile string `yaml:"short_url_file"`
}

var globalConfig *Config

func GetConfig() *Config {
	return globalConfig
}

func Init(filepath string) error {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("read config file=%v failed: %v", filepath, err)
	}

	globalConfig = &Config{}
	err = yaml.Unmarshal(content, globalConfig)
	if err != nil {
		return fmt.Errorf("unmarshal conf%v config: %w", filepath, err)
	}
	return nil
}
