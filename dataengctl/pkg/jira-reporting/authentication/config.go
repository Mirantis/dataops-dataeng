package main

import (
	"io/ioutil"

	"github.com/andygrunwald/go-jira"
	"gopkg.in/yaml.v2"
)

type Config struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

func readConfig() (*Config, error) {
	config := &Config{}
	cfgFile, err := ioutil.ReadFile("~/Git/dataeng/.secrets/jira-config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(cfgFile, config)
	return config, err
}

func getClientFromConfig(config *Config) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Token,
	}
	return jira.NewClient(tp.Client(), config.URL)
}

func getClient() (*jira.Client, error) {
	config, err := readConfig()
	if err != nil {
		return nil, err
	}
	return getClientFromConfig(config)
}
