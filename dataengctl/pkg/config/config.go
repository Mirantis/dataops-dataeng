package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type DataConfig struct {
	// DataConfigs in Current Implementation
	JiraConfig       *JiraConfig       `yaml:"jiraConfig"`
	SalesforceConfig *SalesForceConfig `yaml:"salesForceConfig"`
	GitHubConfig     *GitHubConfig     `yaml:"gitHubConfig"`
	SnowFlakeConfig  *SnowFlakeConfig  `yaml:"snowflakeConfig"`
}

type JiraConfig struct {
	// User Fields Required
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

type SalesForceConfig struct {
	// User Fields Required
	URL        string `yaml:"url"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Token      string `yaml:"token"`
	ClientID   string `yaml:"clientID"`
	APIVersion string `yaml:"apiVersion"`
}

type GitHubConfig struct {
	OrgID string `yaml:"orgID"`
	Token string `yaml:"token"`
}

type SnowFlakeConfig struct {
	// TODO populate this with correct fields
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

func ReadConfig(path string, config *DataConfig) error {
	// Set Config YAML File Path
	cfgFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// Unmarshal YAML FILE
	err = yaml.Unmarshal(cfgFile, config)
	return err
}

// Validate Function for Jira Config
func (c *JiraConfig) ValidateJira() error {
	if c.Token == "" {
		return fmt.Errorf("your jiraConfig is missing a token field")
	}

	if c.URL == "" {
		return fmt.Errorf("your jiraConfig is missing a URL field")
	}

	if c.Username == "" {
		return fmt.Errorf("your jiraConfig is missing a username field")
	}
	return nil
}

// Validate Function for SalesForce Config
func (c *SalesForceConfig) ValidateSalesForce() error {
	if c.Token == "" {
		return fmt.Errorf("your salesforceConfig is missing a token field")
	}

	if c.URL == "" {
		return fmt.Errorf("your salesforceConfig is missing a URL field")
	}

	if c.Username == "" {
		return fmt.Errorf("your salesforceConfig is missing a username field")
	}
	return nil
}

// Validate Function for GitHub Config
func (c *GitHubConfig) ValidateGithub() error {
	if c.Token == "" {
		return fmt.Errorf("your githubConfig is missing a token field")
	}
	if c.OrgID == "" {
		return fmt.Errorf("your githubConfig is missing an orgID field")
	}
	return nil
}

// Validate Function for Snowflake Config
func (c *SnowFlakeConfig) ValidateSnowflake() error {
	if c.Token == "" {
		return fmt.Errorf("your snowflakeConfig is missing a token field")
	}

	if c.Password == "" {
		return fmt.Errorf("your snowflakeConfig is missing a password field")
	}

	if c.Username == "" {
		return fmt.Errorf("your snowflakeConfig is missing a username field")
	}
	return nil
}
