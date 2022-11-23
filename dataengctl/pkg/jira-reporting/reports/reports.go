package reports

import (
	"encoding/json"
	"io"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/andygrunwald/go-jira"
)

type Options struct {
	Output string
	Config client.DataClientInterface
}

func (lo *Options) List(w io.Writer) error {
	jiraClient, err := lo.Config.JiraClient()
	if err != nil {
		return err
	}
	boardList, _, err := jiraClient.GetAllBoards(&jira.BoardListOptions{})
	if err != nil {
		return err
	}
	b, err := json.Marshal(boardList)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}
