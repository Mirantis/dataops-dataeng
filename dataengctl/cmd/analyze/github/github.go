package github

import (
	"os"
	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/github/reports"

	"github.com/spf13/cobra"
)

func NewAnalyzeGitHubCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "github",
		Short:   "analyze github",
		Example: "dataengctl analyze github",
	}
	cmd.AddCommand(NewAnalyzeGitHubCommitHistoryReportCommand(dataClient))
	cmd.AddCommand(NewAnalyzeExpectedContributionsReportCommand(dataClient))
	return cmd
}

//RunExpectedContributionsReport
func NewAnalyzeExpectedContributionsReportCommand(dataClient *client.DataClient) *cobra.Command {
	var path, team string
	cmd := &cobra.Command{
		Use:   "commitaverage",
		Short: "analyze commitaverage",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return reports.RunExpectedContributionsReport(path, team, dataClient, os.Stdout)
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "path for the config file being passed")
	cmd.Flags().StringVarP(&team, "team", "t", "", "give me a team that is defined in the config path, example mke_users")
	return cmd
}

func NewAnalyzeGitHubCommitHistoryReportCommand(dataClient *client.DataClient) *cobra.Command {
	var path, team string
	cmd := &cobra.Command{
		Use:   "commits",
		Short: "analyze commits",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return reports.RunReport(path, team, dataClient, os.Stdout)
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "path for the config file being passed")
	cmd.Flags().StringVarP(&team, "team", "t", "", "give me a team that is defined in the config path, example mke_users")
	return cmd
}
