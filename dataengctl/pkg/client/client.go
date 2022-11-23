package client

import (
	"fmt"

	dataconfig "github.com/Mirantis/dataeng/dataengctl/pkg/config"
	github "github.com/Mirantis/dataeng/dataengctl/pkg/github/client"
	jira "github.com/Mirantis/dataeng/dataengctl/pkg/jira/client"
	salesforce "github.com/Mirantis/dataeng/dataengctl/pkg/salesforce/client"
	snowflake "github.com/Mirantis/dataeng/dataengctl/pkg/snowflake/client"
)

func NewDataClient(configPath string) *DataClient {
	return &DataClient{
		ConfigFilePath: configPath,
	}
}

type DataClient struct {
	jiraClient       *jira.Client
	salesForceClient *salesforce.Client
	gitHubClient     *github.Client
	SnowFlakeClient  *snowflake.Client

	ConfigFilePath string
	k8sConfig      string
	k8sConfigMap   string
	K8sSecret      string
}

type DataClientInterface interface {
	JiraClient() (jira.ClientInterface, error)
	SalesForceClient() (*salesforce.Client, error)
	//GitHubClient() (*github.Client, error)
	//SnowFlakeClient() (*snowflake.Client, error)
}

type KubernetesClient interface {
	KubernetesClient()
}

// func (dc *DataClient) GitHubClient() (*github.Client, error) {
// 	if dc.GitHubClient()!= nil {
// 		return dc.gitHubClient, nil
// 	}

// 	dataConf := &dataconfig.DataConfig{}
// 	err := dataconfig.ReadConfig(dc.ConfigFilePath, dataConf)
// 	if err != nil {
// 		return nil, err
// 	}
// }

func (dc *DataClient) SalesForceClient() (*salesforce.Client, error) {

	if dc.salesForceClient != nil {
		return dc.salesForceClient, nil
	}

	dataConf := &dataconfig.DataConfig{}
	err := dataconfig.ReadConfig(dc.ConfigFilePath, dataConf)
	if err != nil {
		return nil, err
	}

	if dataConf.SalesforceConfig.Token == "" {
		return nil, fmt.Errorf("your salesforceConfig config at '%s' is missing a token field", dc.ConfigFilePath)
	}

	if dataConf.SalesforceConfig.URL == "" {
		return nil, fmt.Errorf("your salesforceConfig config at '%s' is missing a URL field", dc.ConfigFilePath)
	}

	if dataConf.JiraConfig.Username == "" {
		return nil, fmt.Errorf("your salesforceConfig config at '%s' is missing a username field", dc.ConfigFilePath)
	}

	dc.salesForceClient = &salesforce.Client{
		Config: dataConf.SalesforceConfig,
	}

	return dc.salesForceClient, nil
}

func (dc *DataClient) JiraClient() (jira.ClientInterface, error) {
	if dc.jiraClient != nil {
		return dc.jiraClient, nil
	}

	dataConf := &dataconfig.DataConfig{
		JiraConfig:       &dataconfig.JiraConfig{},
		SalesforceConfig: &dataconfig.SalesForceConfig{},
		GitHubConfig:     &dataconfig.GitHubConfig{},
	}
	err := dataconfig.ReadConfig(dc.ConfigFilePath, dataConf)
	if err != nil {
		return nil, err
	}

	if err = dataConf.JiraConfig.ValidateJira(); err != nil {
		return nil, fmt.Errorf("received an error reading from jira config file at '%s', error is \n '%s'",
			dc.ConfigFilePath,
			err.Error())
	}

	dc.jiraClient = &jira.Client{
		Config: dataConf.JiraConfig,
	}

	return dc.jiraClient, nil
}

func (dc *DataClient) GitHubClient() (*github.Client, error) {
	if dc.gitHubClient != nil {
		return dc.gitHubClient, nil
	}

	dataConf := &dataconfig.DataConfig{
		JiraConfig:       &dataconfig.JiraConfig{},
		SalesforceConfig: &dataconfig.SalesForceConfig{},
		GitHubConfig:     &dataconfig.GitHubConfig{},
	}

	err := dataconfig.ReadConfig(dc.ConfigFilePath, dataConf)
	if err != nil {
		return nil, err
	}
	dc.gitHubClient = &github.Client{
		Config: dataConf.GitHubConfig,
	}
	return dc.gitHubClient, nil
}
