package jira

import (
	"os"

	"github.com/Mirantis/dataeng/dataengctl/pkg/analyze/jira/summary"
	"github.com/Mirantis/dataeng/dataengctl/pkg/client"

	"github.com/spf13/cobra"
)

func NewAnalyzeJiraCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "jira",
		Short:   "analyze jira",
		Example: "dataengctl analyze jira",
	}
	cmd.AddCommand(NewAnalyzeJiraIssuesCommand(dataClient))
	return cmd
}

func NewAnalyzeJiraIssuesCommand(dataClient *client.DataClient) *cobra.Command {
	opts := summary.Options{
		DataClient: dataClient,
	}
	cmd := &cobra.Command{
		Use:   "issues",
		Short: "Analyze issues",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.RunE(os.Stdout)
		},
	}
	cmd.Flags().StringVarP(&opts.IssueType, "issue-type", "t", "", "give me type of the issue, i know only Epic")
	cmd.Flags().StringVarP(&opts.ProjectKey, "project-key", "p", "", "give me project key, for example CDC")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "json", "use this flag to specify output")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "use this flag to specify the start-date using yyyy-mm-dd")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "use this flag to specify the end-date using yyyy-mm-dd")
	cmd.Flags().BoolVarP(&opts.Summary, "summary", "s", false, "use this flag to specify summary report")
	cmd.Flags().IntVar(&opts.TimeoutSeconds, "timeout", 120, "timeout specified in seconds")

	return cmd
}
