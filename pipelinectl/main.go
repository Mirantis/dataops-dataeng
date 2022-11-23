package main

import (
	"fmt"
	"os"

	"github.com/Mirantis/dataeng/pipelinectl/cmd"
)

func main() {
	if err := cmd.NewPipelineCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
