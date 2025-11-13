package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var currentVersion string = ""

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints dPanel ID CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(currentVersion)
		},
	}
}
