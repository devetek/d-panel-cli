package main

import (
	"fmt"

	"github.com/devetek/d-panel-cli/internal/api"
	"github.com/devetek/d-panel-cli/internal/helper"
	"github.com/devetek/d-panel-cli/internal/logger"
	"github.com/devetek/d-panel-cli/internal/tunnel"
	"github.com/devetek/tuman/pkg/marijan"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type TunnelCmd struct {
	cmd       *cobra.Command
	zapLogger *zap.Logger

	tunnelHttpListener string
	tunnelHttpService  string
	tunnelSshListener  string
	tunnelSshService   string
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
			// init dPanel client
			client := api.NewClient()

			// check if session exist
			err := client.CheckSessionExist()
			if err != nil {
				logger.Error("Please login to your dPanel account, use command 'dnocs auth login --email=\"email@email.com\" --password=\"password\"'")
				return
			}

			// TODO: Remove sync communication after MQTT architecture completed!
			if m.tunnelHttpListener == "" {
				logger.Error("Please set HTTP public listener in the tunnel")
				return
			}

			if m.tunnelSshListener == "" {
				logger.Error("Please set SSH public listener in the tunnel")
				return
			}

			if !helper.IsSudo() {
				logger.Error("You must run this command as sudo, currenty tunnel required to running under root")
				return
			}

			if helper.IsPortUsed(tunnel.TunnelHost, m.tunnelHttpListener) {
				logger.Error("Port already in used in the tunnel server, choose another HTTP port or contact prakasa@devetek.com")
				return
			}

			if helper.IsPortUsed(tunnel.TunnelHost, m.tunnelSshListener) {
				logger.Error("Port already in used in the tunnel server, choose another SSH port or contact prakasa@devetek.com")
				return
			}

			var tunnelCreation = tunnel.NewTunnel().SetConfig([]marijan.Config{
				{
					NoTCP:        false,
					ID:           fmt.Sprintf("ssh-%s-to-%s", m.tunnelSshListener, m.tunnelSshService),
					TunnelHost:   tunnel.TunnelHost,
					TunnelPort:   tunnel.TunnelPort,
					ListenerHost: "0.0.0.0",
					ListenerPort: m.tunnelSshListener,
					ServiceHost:  "localhost",
					ServicePort:  m.tunnelSshService,
					State:        marijan.ConfigStateActive,
				},
				{
					NoTCP:        false,
					ID:           fmt.Sprintf("http-%s-to-%s", m.tunnelHttpListener, m.tunnelHttpService),
					TunnelHost:   tunnel.TunnelHost,
					TunnelPort:   tunnel.TunnelPort,
					ListenerHost: "0.0.0.0",
					ListenerPort: m.tunnelHttpListener,
					ServiceHost:  "localhost",
					ServicePort:  m.tunnelHttpService,
					State:        marijan.ConfigStateActive,
				},
			})

			err = tunnelCreation.Download()
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

	runCmd.PersistentFlags().StringVarP(&m.tunnelHttpListener, "tunnel-http-listener", "", "", "Public SSH listener to your machine")
	runCmd.PersistentFlags().StringVarP(&m.tunnelHttpService, "tunnel-http-service", "", "9000", "HTTP port of your machine")
	runCmd.PersistentFlags().StringVarP(&m.tunnelSshListener, "tunnel-ssh-listener", "", "", "Public SSH listener to your machine")
	runCmd.PersistentFlags().StringVarP(&m.tunnelSshService, "tunnel-ssh-service", "", "22", "SSH port of your machine")

	return runCmd
}
