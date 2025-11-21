package main

import (
	"github.com/devetek/d-panel-cli/internal/logger"
	"github.com/devetek/d-panel-cli/internal/tunnel"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type TunnelCmd struct {
	cmd       *cobra.Command
	zapLogger *zap.Logger

	// sshIP    string
	// sshPort  string
	// httpPort string
}

func NewTunnelCmd(logger *zap.Logger) *TunnelCmd {
	return &TunnelCmd{
		zapLogger: logger,
		cmd: &cobra.Command{
			Use:   "tunnel",
			Short: "Manage dPanel tunnel",
		},
	}
}

func (m *TunnelCmd) Connect() *cobra.Command {
	m.cmd.AddCommand(
		m.create(),
	)

	return m.cmd
}

func (m *TunnelCmd) create() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "create",
		Short: "Open connection to tunnel",
		Long:  `Create public access to this machine use tunnel, make it accessible from dPanel.`,
		Run: func(cmd *cobra.Command, args []string) {
			// if !helper.IsSudo() {
			// 	logger.Error("You must run this command as sudo")
			// 	return
			// }

			// // init dPanel client
			// client := api.NewClient()

			// // check if session exist
			// err := client.CheckSessionExist()
			// if err != nil {
			// 	logger.Error("Error check session exist: " + err.Error())
			// 	return
			// }

			var tunnelCreation = tunnel.NewTunnel()

			err := tunnelCreation.Download()
			if err != nil {
				logger.Error(err.Error())
			}

			err = tunnelCreation.Extract()
			if err != nil {
				logger.Error(err.Error())
			}

			err = tunnelCreation.CreateService()
			if err != nil {
				logger.Error(err.Error())
			}

		},
	}

	// runCmd.PersistentFlags().StringVarP(&m.sshIP, "ssh-ip", "i", "", "SSH IP of your machine")
	// runCmd.PersistentFlags().StringVarP(&m.sshPort, "ssh-port", "s", "22", "SSH port of your machine")
	// runCmd.PersistentFlags().StringVarP(&m.httpPort, "http-port", "t", "9000", "HTTP port of your machine")

	return runCmd
}
