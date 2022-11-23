package cmd

import (
	//"os"

	"github.com/spf13/cobra"
)

func NewPipelineCommand() *cobra.Command {
	var debug bool
	var configPath string
	// dataClient := &client.DataClient{}
	cmd := &cobra.Command{
		Use:   "pipelinectl",
		Short: "creates pipelines to and from various data sources",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// log.Init(debug, os.Stderr)
			// *dataClient = *client.NewDataClient(configPath)
		},
	}
	// Wrapper for Output Flag to All Base Commands
	/*func NewOutputCommand() *cobra.Command {
		var output string
		cmd := &cobra.Command{
			Use:	"output",
			Short: 	"output data format to csv, dataframe, json, yaml (DEFAULT is JSON)",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				log.Init(debug, os.Stderr)
		},
	}*/
	// Persistent Global Flags
	cmd.PersistentFlags().StringVar(&configPath, "config-file", "", "path to config file")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "specify debug level")
	//cmd.PersistentFlags().StringVar(&output), "output", false, "specify output format on data return")

	// Persistent Global Commands
	// cmd.AddCommand(analyze.NewAnalyzeCmd(dataClient))
	// cmd.AddCommand(jira.NewJiraCmd(dataClient))
	// cmd.AddCommand(salesforce.NewSalesforceCmd(dataClient))
	return cmd
}