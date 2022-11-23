package main

import (
	"fmt"
	"os"

	"github.com/Mirantis/dataeng/dataengctl/cmd"
)

func main() {
	if err := cmd.NewDataengCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
