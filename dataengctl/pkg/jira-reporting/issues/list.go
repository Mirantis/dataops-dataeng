package issues

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"strings"

	dataclient "github.com/Mirantis/dataeng/dataengctl/pkg/client"
	jiraclient "github.com/Mirantis/dataeng/dataengctl/pkg/jira/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/log"
)

const (
	// TODO Make all const configurable
	// baseQuery            = `issuetype = "%s" AND project = "%s" ORDER BY createdDate ASC`
	orderByDateSortQuery = `ORDER BY createdDate ASC`
	issueTypeQuery       = `issuetype = "%s"`
	projectKeyQuery      = `project = "%s"`
	orderQuery           = `ORDER BY createdDate %s`
)

type Options struct {
	IssueType  string
	ProjectKey string
	Output     string
	DataClient dataclient.DataClientInterface
}

func (o *Options) List(w io.Writer) error {
	jiraClient, err := o.DataClient.JiraClient()
	if err != nil {
		return err
	}
	query, err := o.getQuery()
	if err != nil {
		return err
	}

	iterator, err := jiraClient.IssueIterator(jiraclient.IssueIteratorOptions{
		PageSize: 10,
		Query:    query,
		Expand:   "changelog",
	})
	if err != nil {
		return err
	}

	// TODO add timeout to the context and cli flag
	for {
		issue, err := iterator.Next(context.Background())
		if err != nil {
			return err
		}
		if issue == nil {
			break // if issue is nil and there is no error that means we got no more issues.
		}

		issue.Fields.Unknowns = nil
		b, err := json.Marshal(issue)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		if err != nil {
			return err
		}

	}
	return nil
}

func (lo *Options) getQuery() (string, error) {
	var finalQuery string

	if lo.IssueType == "" {
		return "", fmt.Errorf("issue type must provided")
	}

	var typeQuery = fmt.Sprintf(issueTypeQuery, lo.IssueType)
	finalQuery = typeQuery

	if lo.ProjectKey != "" {
		projectQuery := fmt.Sprintf(projectKeyQuery, lo.ProjectKey)
		finalQuery = strings.Join([]string{typeQuery, projectQuery}, " AND ")
	}

	// TODO (rbarrett) Make ordering configurable
	finalQuery = strings.Join([]string{finalQuery, orderByDateSortQuery}, " ")

	log.Printf("Final query '%s'", finalQuery)
	return finalQuery, nil
}
