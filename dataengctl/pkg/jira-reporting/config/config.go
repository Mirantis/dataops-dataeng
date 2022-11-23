package config

import (
	"os"

	jira "github.com/andygrunwald/go-jira"
	"gopkg.in/yaml.v2"
)

type Config struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

func ReadConfig(path string, config *Config) error {
	cfgFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(cfgFile, config)
	return err
}

func (c *Config) Client() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: c.Username,
		Password: c.Token,
	}
	return jira.NewClient(tp.Client(), c.URL)
}
