package summary

import (
	"sort"
	"time"

	"github.com/andygrunwald/go-jira"
)

const (
	StatusBlocked                 = "Blocked"
	StatusBacklog                 = "Backlog"
	StatusDone                    = "Done"
	StatusClosed                  = "Closed"
	StatusReleased                = "Released"
	StatusNew                     = "New"
	StatusToDo                    = "To Do"
	StatusInProgress              = "In Progress"
	StatusInReview                = "In Review"
	StatusInValidation            = "In Validation"
	StatusPendingRelease          = "Pending Release"
	StatusSustainingTriage        = "sustaining triage"
	StatusSustainingInvestigating = "Sustaining: Investigating"
	StatusInvestigating           = "investigating"
	StatusSustainingInProgress    = "Sustaining: In Progress"
	StatusDevTriage               = "Dev Triage"

	FieldNameStatus = "status"
)

// TODO Add FIelds for custom field to get Customer's Name, AccountID, the URL
// TODO Give all FIELD tickets that are open for renewall within next quarter
// TODO Make Fields Configurable
// TODO Filter off --start-date --end-date
// TODO Configurable Maps for Team Status and Breakdown Structures

type IssueReport struct {
	IssueKey             string
	Priority             string
	Assignee             string
	TimesStatusChanged   int
	TicketDuration       time.Duration
	InitialTicketStatus  string
	EngineersTouched     map[string]struct{} //TODO
	CurrentStatus        string
	TicketReporter       string
	TimeInStatus         map[string]time.Duration
	AverageTimePerStatus string
}

// TODO
// REFERENCE Struct for https://mirantis.jira.com/browse/DATAENG-3
// REVISIT ACCEPTANCE CRITERIA
// ./dataengctl analyze  jira issues --issue-type <issue_key> --project-key <project_key> --summary --config-file ${HOME}/.dataeng/config.yaml  | jq "."

type SummaryStatisticsIssuesReport struct {
	ProjectId                        string                   // Project Key Field
	IssuesCount                      uint64                   // Total Number of Issues within Project
	AverageTicketDuration            time.Duration            // Final State - Ending State
	AverageTicketEngineerTouches     uint64                   // Average Amount of Engineers Touched
	AverageTimeToCompletionPerStatus map[string]time.Duration // Average Time to Complete By Status
	AverageTicketStatusTimes         map[string]time.Duration // Average
	NumberOfIssuesClosed             uint64
	totalTicketEnginerTouches        map[string]int
	totalTimeToCompletionPerStatus   map[string]time.Duration
	totalTicketStatusTimes           map[string]time.Duration
	totalTicketDuration              time.Duration
	totalEngineersTouched            uint64
}

// TODO
// For Every Issue List Fields
// Limit Fields to Specific Set of Fields HARDCODED
// NEEDED FIELD TYPES: id, emailAddress, displayName, active, timezone,
// Append IssueReport to Each Issue in Single JSON Object
// --summary provides summary output for issues listed as SummaryStatisticsIssuesReport

func NewIssueReport(issue *jira.Issue) (IssueReport, error) {
	timeInStatus, err := TimeInStatus(issue)
	if err != nil {
		return IssueReport{}, err
	}

	ticketDuration, touched, err := calculateTicketDurationAndAuthors(issue)
	if err != nil {
		return IssueReport{}, err
	}
	assignee := ""
	if issue.Fields.Assignee != nil {
		assignee = issue.Fields.Assignee.DisplayName
	}

	reporter := ""
	if issue.Fields.Reporter != nil {
		reporter = issue.Fields.Reporter.Name
	}

	priority := ""
	if issue.Fields.Priority != nil {
		priority = issue.Fields.Priority.Name
	}

	currentStatus := ""
	if issue.Fields.Status != nil {
		currentStatus = issue.Fields.Status.Name
	}

	averageTimePerStatus := ""
	if issue.Fields.Status != nil {
		averageTimePerStatus = issue.Fields.Status.Name
	}

	return IssueReport{
		IssueKey:             issue.Key,
		Priority:             priority,
		CurrentStatus:        currentStatus,
		Assignee:             assignee,
		TicketReporter:       reporter,
		TicketDuration:       ticketDuration,
		TimeInStatus:         timeInStatus,
		AverageTimePerStatus: averageTimePerStatus,
		EngineersTouched:     touched,
	}, nil
}

func calculateTicketDurationAndAuthors(issue *jira.Issue) (time.Duration, map[string]struct{}, error) {
	createdAt := time.Time(issue.Fields.Created)
	touched := map[string]struct{}{}

	var endedAt time.Time
	var err error
	for _, historyRecord := range issue.Changelog.Histories {
		touched[historyRecord.Author.DisplayName] = struct{}{}
		for _, item := range historyRecord.Items {
			if item.Field == FieldNameStatus {
				if terminalStatus(item.ToString) {
					endedAt, err = historyRecord.CreatedTime()
					if err != nil {
						return time.Duration(0), touched, err
					}
				}
			}
		}
	}
	if endedAt.IsZero() {
		endedAt = time.Now()
	}
	return endedAt.Sub(createdAt), touched, nil

}

type timeContainer struct {
	name      string
	startDate time.Time
}

type sorted struct {
	s []timeContainer
}

func (s *sorted) Len() int {
	return len(s.s)
}

func (s *sorted) Swap(i, j int) {
	tmp := s.s[j]
	s.s[j] = s.s[i]
	s.s[i] = tmp
}

func (s *sorted) Less(i, j int) bool {
	return s.s[i].startDate.Before(s.s[j].startDate)
}

func TimeInStatus(issue *jira.Issue) (map[string]time.Duration, error) {
	result := map[string]time.Duration{}

	tmpSlice := []timeContainer{}
	for _, historyRecord := range issue.Changelog.Histories {
		for _, item := range historyRecord.Items {
			if item.Field == FieldNameStatus {
				endedAt, err := historyRecord.CreatedTime()
				if err != nil {
					return nil, err
				}
				tmpSlice = append(tmpSlice, timeContainer{
					startDate: endedAt,
					name:      item.ToString,
				})
			}
		}
	}
	sortedSlice := &sorted{s: tmpSlice}
	sort.Sort(sortedSlice)
	for index, statusChange := range sortedSlice.s {
		if index+1 >= len(sortedSlice.s) {
			result[statusChange.name] = 0
			return result, nil
		}
		result[statusChange.name] = sortedSlice.s[index+1].startDate.Sub(statusChange.startDate)
	}
	return nil, nil
}

func terminalStatus(status string) bool {
	switch status {
	case StatusClosed, StatusDone, StatusReleased, StatusPendingRelease:
		return true
	default:
		return false
	}
}

func (s *SummaryStatisticsIssuesReport) Finalize() error {
	if s.IssuesCount != 0 {
		s.AverageTicketDuration = s.totalTicketDuration / time.Duration(s.IssuesCount)
		for status, t := range s.totalTimeToCompletionPerStatus {
			s.AverageTimeToCompletionPerStatus[status] = t / time.Duration(s.IssuesCount)
		}
		s.AverageTicketEngineerTouches = s.totalEngineersTouched / s.IssuesCount
	}

	return nil
}
