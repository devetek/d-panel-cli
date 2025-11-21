package main

import (
	"runtime"

	"github.com/devetek/d-panel-cli/internal/logger"
	"github.com/spf13/cobra"
)

func systemInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Prints the version",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Normal("Your System Information:")
			logger.Success("OS: " + runtime.GOOS)
			logger.Success("Arch: " + runtime.GOARCH)
		},
	}
}
