package issues

import (
	"os"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/jira-reporting/issues"
	"github.com/spf13/cobra"
)

func NewIssueCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Short:   "interact with issues",
		Example: "dataengctl jira issues list",
	}
	cmd.AddCommand(NewListCommand(dataClient))
	return cmd
}

func NewListCommand(dataClient *client.DataClient) *cobra.Command {
	issueListOptions := &issues.Options{
		DataClient: dataClient,
	}
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List jira epics",
		Example:       "",
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return issueListOptions.List(os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&issueListOptions.IssueType, "issue-type", "t", "", "give me type of the issue, i know only Epic")
	cmd.Flags().StringVarP(&issueListOptions.ProjectKey, "project-key", "p", "", "give me project key, for example CDC")

	cmd.Flags().StringVarP(&issueListOptions.Output, "output", "o", "json", "use this flag to specify output")
	return cmd
}
