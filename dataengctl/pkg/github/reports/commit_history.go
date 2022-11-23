package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/shurcooL/graphql"
	"gopkg.in/yaml.v2"
)

func NewReport(weeks []*Week) *Report {
	return &Report{weeks: weeks}
}

type Report struct {
	weeks []*Week
}

func (r *Report) AverageContributionsCount() int {
	weekCount := len(r.weeks)
	totalContributions := 0
	for _, week := range r.weeks {
		totalContributions += week.Contributions()
	}
	return totalContributions / weekCount
}

type Teams struct {
	Map map[string]Team `yaml:"map"`
}

type Team []string

type Week struct {
	StartDate       string         `json:"Week"`
	ContributionMap map[string]int `json:"Contributions"`
}

type SummaryReport struct {
	StartDate                  string `json:"Week"`
	NumberofWeeks              int    `json:"Number of Weeks"`
	ContributionCount          int    `json:" Countrbution Count"`
	ContributionGrowthOverTime int    `json:"Contribution Growth Over Time"`
	ContributionUtility        int    `json:"Contribution Utility"`
	ContributionComplexity     int    `json:"Contribution Complexity"`
	ContributionsAddedCount    int    `json:"Contributions Added Count"`
	ContributionsDeletedCount  int    `json:"Contributions Deleted Count"`
	ExpectedContributions      int    `json:"Expected Contributions"`
	ExternalContributions      int    `json:"External Contributions"`
	GitHubIssuesClosed         int    `json:"GitHub Issues Closed"`
	GitHubIssuesOpened         int    `json:"GitHub Issues Opened"`
	InternalContributions      int    `json:"Internal Contributions"`
	OrganizationsImpacted      string `json:"Organizations Impacted"`
	PersonalContributions      int    `json:"Personal Contributions"`
	PullRequestsClosedCount    int    `json:"Pull Requests Closed"`
	PullRequestsOpenedCount    int    `json:"Pull Reuqests Opened"`
	TeamContributions          int    `json:" Team Contributions"`
}
type AverageSummaryReport struct {
	StartDate                         string `json:"Week"`
	NumberofWeeks                     int    `json:"Number of Weeks"`
	AverageContributionCount          int    `json:"Average Countrbution Count"`
	AverageContributionGrowthOverTime int    `json:"Contribution Growth Over Time"`
	AverageContributionUtility        int    `json:"Contribution Utility"`
	AverageContributionComplexity     int    `json:"Contribution Complexity"`
	AverageContributionsAdded         int    `json:"Average Contributions Added"`
	AverageContributionsDeleted       int    `json:"Average Contributions Deleted"`
	AverageExpectedContributions      int    `json:"Average Expected Contributions"`
	AverageExternalContributions      int    `json:"Average External Contributions"`
	AverageGitHubIssuesClosed         int    `json:"Average GitHub Issues Closed"`
	AverageGitHubIssuesOpened         int    `json:"Average GitHub Issues Opened"`
	AverageInternalContributions      int    `json:"Average Internal Contributions"`
	AverageOrganizationsImpacted      string `json:"Average Organizations Impacted"`
	AveragePersonalContributions      int    `json:"Average Personal Contributions"`
	AveragePullRequestsClosedCount    int    `json:"Average Pull Requests Closed"`
	AveragePullRequestsOpenedCount    int    `json:"Average Pull Reuqests Opened"`
	AverageTeamContributions          int    `json:"Average Team Contributions"`
}

func getTeams(path string) (*Teams, error) {
	teams := &Teams{}
	cfgFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return teams, yaml.Unmarshal(cfgFile, teams)
}

var apiURL = "https://api.github.com/graphql"
var endTime = time.Now().UTC()
var startTime = endTime.Add(-365 * 24 * time.Hour)

type DateTime string
type Contributions struct {
	User struct {
		Name                    graphql.String
		ContributionsCollection struct {
			ContributionCalendar struct {
				Weeks []struct {
					ContributionDays []struct {
						ContributionCount graphql.Int
						Date              graphql.String
					}
					FirstDay graphql.String
				}
			}
		} `graphql:"contributionsCollection(organizationID: $orgID from: $startTime to: $endTime)"`
	} `graphql:"user(login: $login)"`
}

type TokenAuthTransport struct {
	Token string
}

func (t TokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println(t.Token)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", t.Token))
	return http.DefaultTransport.RoundTrip(req)
}

func RunReport(path, team string, dataClient *client.DataClient, w io.Writer) error {
	weeks, err := GetWorkWeeks(path, team, dataClient)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(weeks)
}

// TODO FIX FLOAT64 TYPE FOR AVERAGES
func RunExpectedContributionsReport(path, team string, dataClient *client.DataClient, w io.Writer) error {
	workDone, err := ExpectedContributions(path, team, dataClient)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(workDone)
}

func GetWorkWeeks(path, team string, dataClient *client.DataClient) ([]*Week, error) {
	teams, err := getTeams(path)
	if err != nil {
		return nil, err
	}
	users, exist := teams.Map[team]
	if !exist {
		return nil, fmt.Errorf("team '%s' doesn't exist in the config file at %s", team, path)
	}

	gitHubClient, err := dataClient.GitHubClient()
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]int)
	for _, user := range users {
		result[user] = make(map[string]int)
		client := graphql.NewClient(apiURL, &http.Client{
			Transport: TokenAuthTransport{
				Token: gitHubClient.Config.Token,
			},
		})

		c := &Contributions{}
		vars := map[string]interface{}{
			"login":     graphql.String(user),
			"orgID":     graphql.ID(gitHubClient.Config.OrgID),
			"startTime": DateTime(startTime.Format(time.RFC3339)),
			"endTime":   DateTime(endTime.Format(time.RFC3339)),
		}
		if err := client.Query(context.Background(), c, vars); err != nil {
			fmt.Printf("query error: %v\n", err)
			continue
		}

		for _, week := range c.User.ContributionsCollection.ContributionCalendar.Weeks {
			weeklyCount := 0
			for _, day := range week.ContributionDays {
				weeklyCount += int(day.ContributionCount)
			}
			result[user][string(week.FirstDay)] = weeklyCount
		}
	}

	weeks := make([]string, 0)
	for week, _ := range result[users[0]] {
		weeks = append(weeks, week)
	}
	sort.Strings(weeks)

	returnWeeks := []*Week{}

	for _, week := range weeks {
		weekJson := Week{
			StartDate:       week,
			ContributionMap: map[string]int{},
		}
		for _, user := range users {
			weekJson.ContributionMap[user] = result[user][week]
		}
		returnWeeks = append(returnWeeks, &weekJson)
	}

	return returnWeeks, nil
}

func (w *Week) ListContributors() []string {
	contributors := []string{}
	for key, _ := range w.ContributionMap {
		contributors = append(contributors, key)
	}
	return contributors
}

func (w *Week) Contributions() int {
	totalContributions := 0
	for _, contributions := range w.ContributionMap {
		totalContributions += contributions
	}
	return totalContributions
}

func (w *Week) AverageContributions() float64 {
	averageContributions := 0
	for _, contributions := range w.ContributionMap {
		averageContributions += (contributions)/5
	}
	return float64(averageContributions)
}

func (w *Week) NonContributor() bool {
	//if Contributions := 0
	//for _, week := range w.ConGetWorkWeeks()
	return false
}

func RunSummaryReport(path, team string, dataClient *client.DataClient) error {
	return nil
}

func StartDate(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func NumberofWeeks(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ContributionCount(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ContributionCalendar(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ContributionComplexity(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ContributionDays(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ContributionGrowthOverTime(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ContributionUtility(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ContributionsAddedCount(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ContributionsDeletedCount(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func ExpectedContributions(path, team string, dataClient *client.DataClient) ([]*Week, error) {
	teams, err := getTeams(path)
	if err != nil {
		return nil, err
	}
	users, exist := teams.Map[team]
	if !exist {
		return nil, fmt.Errorf("team '%s' doesn't exist in the config file at %s", team, path)
	}

	gitHubClient, err := dataClient.GitHubClient()
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]int)
	for _, user := range users {
		result[user] = make(map[string]int)
		client := graphql.NewClient(apiURL, &http.Client{
			Transport: TokenAuthTransport{
				Token: gitHubClient.Config.Token,
			},
		})

		c := &Contributions{}
		vars := map[string]interface{}{
			"login":     graphql.String(user),
			"orgID":     graphql.ID(gitHubClient.Config.OrgID),
			"startTime": DateTime(startTime.Format(time.RFC3339)),
			"endTime":   DateTime(endTime.Format(time.RFC3339)),
		}
		if err := client.Query(context.Background(), c, vars); err != nil {
			fmt.Printf("query error: %v\n", err)
			continue
		}

		for _, weeklyaverage := range c.User.ContributionsCollection.ContributionCalendar.Weeks {
			averageCount := 0
			for _, day := range weeklyaverage.ContributionDays {
				averageCount += int(day.ContributionCount)/5
				//float64(day.ContributionCount)/7
			}
			result[user][string(weeklyaverage.FirstDay)] = averageCount
		}
	}

	averageoverweeks := make([]string, 0)
	for weeklyaverage, _ := range result[users[0]] {
		averageoverweeks = append(averageoverweeks, weeklyaverage)
	}
	sort.Strings(averageoverweeks)

	returnWeekkyAverage := []*Week{}

	for _, week := range averageoverweeks {
		weekJson := Week{
			StartDate:       week,
			ContributionMap: map[string]int{},
		}
		for _, user := range users {
			weekJson.ContributionMap[user] = result[user][week]
		}
		returnWeekkyAverage = append(returnWeekkyAverage, &weekJson)
	}

	return returnWeekkyAverage, nil
}

func ExternalContributions(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func GitHubIssuesClosed(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func GitHubIssuesOpened(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func InternalContributions(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func OrganizationsImpactedfunc(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func PersonalContributions(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func PullRequestsClosedCount(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func PullRequestsOpenedCount(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func TeamContributions(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}

func GetContributionUsers(path, team string, dataClient *client.DataClient, w io.Writer) error {
	return nil
}