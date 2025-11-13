package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dpid",
	Short: "dPanel ID CLI",
	Long: `
dPanel ID CLI is a simple CLI tool to interact with dPanel, simplify the process of managing your dPanel.

Full documentation is available at: https://cloud.terpusat.com/docs/d-panel-cli
`,
}

func init() {
	rootCmd.AddCommand(
		versionCmd(),
		runCmd(),
	)
}

func Execute() {
	rootCmd.Version = currentVersion
	cobra.CheckErr(rootCmd.Execute())
}

func main() {
	Execute()
}
