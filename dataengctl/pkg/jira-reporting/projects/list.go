package projects

import (
	"fmt"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
)

func List(dataClient client.DataClientInterface) error {
	jiraClient, err := dataClient.JiraClient()
	if err != nil {
		return err
	}
	body, err := jiraClient.ProjectBytes()
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
