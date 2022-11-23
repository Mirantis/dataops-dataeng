package jira

import (
	"github.com/Mirantis/dataeng/dataengctl/cmd/jira/boards"
	"github.com/Mirantis/dataeng/dataengctl/cmd/jira/dashboards"
	"github.com/Mirantis/dataeng/dataengctl/cmd/jira/issues"
	"github.com/Mirantis/dataeng/dataengctl/cmd/jira/projects"
	"github.com/Mirantis/dataeng/dataengctl/cmd/jira/reports"
	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/spf13/cobra"
)

// Wrapper for Jira Client
func NewJiraCmd(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jira",
		Short: "Interact with jira",
	}

	// Wrapper for Output Flag to All Base Commands
	/*func OutPutCmd() *cobra.Command {
		var output string
		output := &output{}
		cmd := &cobra.Command{
			Use:	"output"
		},
	}*/

	//cmd.PersistentFlags().StringVar(&output), "output", "", "csv,dataframe,json,yaml")

	// Universal Base Commands
	cmd.AddCommand(boards.NewBoardsCommand(dataClient))
	cmd.AddCommand(dashboards.NewDashboardCommand(dataClient))
	cmd.AddCommand(issues.NewIssueCommand(dataClient))
	cmd.AddCommand(projects.NewProjectCommand(dataClient))
	cmd.AddCommand(reports.NewReportsCommand(dataClient))
	//cmd.AddCommand(epics.NewEpicCommand())
	return cmd
}
