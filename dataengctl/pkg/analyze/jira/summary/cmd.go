package summary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	dataclient "github.com/Mirantis/dataeng/dataengctl/pkg/client"
	jiraclient "github.com/Mirantis/dataeng/dataengctl/pkg/jira/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/log"
	"github.com/andygrunwald/go-jira"
)

const (
	// TODO Make all const configurable
	// baseQuery            = `issuetype = "%s" AND project = "%s" ORDER BY createdDate ASC`
	orderByDateSortQuery = `ORDER BY createdDate ASC`
	issueTypeQuery       = `issuetype = "%s"`
	projectKeyQuery      = `project = "%s"`
	dateQuery            = `created %s %s`
	orderQuery           = `ORDER BY createdDate %s`
)

type Options struct {
	IssueType  string
	ProjectKey string
	Output     string
	StartDate  string
	EndDate    string

	Summary        bool
	TimeoutSeconds int
	DataClient     dataclient.DataClientInterface
	summaryReport  *SummaryStatisticsIssuesReport
}

func (o *Options) RunE(w io.Writer) error {
	if o.TimeoutSeconds == 0 {
		o.TimeoutSeconds = 120
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.TimeoutSeconds))
	defer cancel()
	if o.Summary {
		return o.Analyze(ctx, w)
	}
	return o.List(ctx, w)
}

func (o *Options) List(ctx context.Context, w io.Writer) error {
	jiraClient, err := o.DataClient.JiraClient()
	if err != nil {
		return err
	}
	query, err := o.getQuery()
	if err != nil {
		return err
	}

	iterator, err := jiraClient.IssueIterator(jiraclient.IssueIteratorOptions{
		PageSize: 100,
		Query:    query,
		Expand:   "changelog",
	})

	// TODO add timeout to the context and cli flag
	for {
		issue, err := iterator.Next(ctx)
		if err != nil {
			return err
		}
		if issue == nil {
			break // if issue is nil and there is no error that means we got no more issues.
		}
		report, err := NewIssueReport(issue)
		if err != nil {
			return err
		}
		b, err := json.Marshal(report)
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

func (o *Options) Analyze(ctx context.Context, w io.Writer) error {
	jiraClient, err := o.DataClient.JiraClient()
	if err != nil {
		return err
	}
	query, err := o.getQuery()
	if err != nil {
		return err
	}

	iterator, err := jiraClient.IssueIterator(jiraclient.IssueIteratorOptions{
		PageSize: 100,
		Query:    query,
		Expand:   "changelog",
	})
	if err != nil {
		return err
	}

	o.summaryReport = &SummaryStatisticsIssuesReport{
		ProjectId:                        o.ProjectKey,
		totalTimeToCompletionPerStatus:   make(map[string]time.Duration),
		totalTicketEnginerTouches:        make(map[string]int),
		totalTicketStatusTimes:           make(map[string]time.Duration),
		AverageTimeToCompletionPerStatus: make(map[string]time.Duration),
		AverageTicketStatusTimes:         make(map[string]time.Duration),
	}

	// TODO add timeout to the context and cli flag
	for {
		issue, err := iterator.Next(ctx)
		if err != nil {
			return err
		}
		if issue == nil {
			break // if issue is nil and there is no error that means we got no more issues.
		}
		report, err := NewIssueReport(issue)
		if err != nil {
			return err
		}

		if err = o.summaryReport.AddIssue(report, issue); err != nil {
			return err
		}

	}
	err = o.summaryReport.Finalize()
	if err != nil {
		return err
	}

	b, err := json.Marshal(o.summaryReport)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (s *SummaryStatisticsIssuesReport) AddIssue(report IssueReport, issue *jira.Issue) error {
	s.IssuesCount++
	if report.CurrentStatus == "Closed" || report.CurrentStatus == "Done" {
		s.NumberOfIssuesClosed++
	}
	s.totalEngineersTouched += uint64(len(report.EngineersTouched))

	s.totalTicketDuration += report.TicketDuration
	for status, t := range report.TimeInStatus {
		_, exists := s.totalTimeToCompletionPerStatus[status]
		if exists {
			s.totalTimeToCompletionPerStatus[status] += t
		} else {
			s.totalTimeToCompletionPerStatus[status] = t
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

	if lo.StartDate != "" {
		startDateQuery := fmt.Sprintf(dateQuery, ">=", lo.StartDate)
		finalQuery = strings.Join([]string{finalQuery, startDateQuery}, " AND ")
	}

	if lo.EndDate != "" {
		endDateQuery := fmt.Sprintf(dateQuery, "<=", lo.EndDate)
		finalQuery = strings.Join([]string{finalQuery, endDateQuery}, " AND ")
	}

	// TODO (rbarrett) Make ordering configurable
	finalQuery = strings.Join([]string{finalQuery, orderByDateSortQuery}, " ")

	log.Printf("Final query '%s'", finalQuery)
	return finalQuery, nil
}
