package projects

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mirantis/dataeng/dataengctl/pkg/client"
	"github.com/Mirantis/dataeng/dataengctl/pkg/jira-reporting/projects"
)

const (
	ListExample = `
# this will print list of the projects
dataengctl jira project list

# this will create output in a yaml format
dataengctl jira project list --output yaml

`
)

func NewProjectCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "interact with projects",
	}

	cmd.AddCommand(NewListCommand(dataClient))
	cmd.AddCommand(NewExpandCommand(dataClient))
	return cmd
}

func NewListCommand(dataClient *client.DataClient) *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:           "list",
		Short:         "List jira projects",
		Example:       ListExample,
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return projects.List(dataClient)
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "yaml", "use this flag to specify output")
	return cmd
}

func NewExpandCommand(dataClient *client.DataClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "expand",
		Short:   "expand jira projects",
		Example: "Some example how to run this",

		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("here is the list of expanded projects\n")
			return nil
		},
	}
	return cmd
}
