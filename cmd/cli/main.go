package main

import (
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "dnocs",
	Short: "dnocs ID CLI",
	Long: `
Simplify the process of managing resource such as user, machine, and application in dPanel (Devetek Panel).

Full documentation is available at: https://cloud.terpusat.com/docs/
`,
}

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logger.Sync()

	rootCmd.AddCommand(
		NewAuthCmd(logger).Connect(),
		NewTunnelCmd(logger).Connect(),
		NewMachineCmd(logger).Connect(),
		versionCmd(),
		systemInfoCmd(),
	)
}

func Execute() {
	rootCmd.Version = currentVersion
	cobra.CheckErr(rootCmd.Execute())
}

func main() {
	Execute()
}
