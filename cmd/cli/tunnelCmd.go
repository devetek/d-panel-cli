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
	"golang.org/x/mod/semver"
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
		m.upgrade(),
		m.create(),
	)

	return m.cmd
}

func (m *TunnelCmd) upgrade() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade marijan binary",
		Long:  fmt.Sprintf("Upgrade marijan binary to the latest version, check latest version in %s", tunnel.BinaryBaseURL),
		Run: func(cmd *cobra.Command, args []string) {
			// init dPanel client
			client := api.NewClient()

			// check if session exist
			err := client.CheckSessionExist()
			if err != nil {
				logger.Error("Please login to your dPanel account, use command 'dnocs auth login --email=\"email@email.com\" --password=\"password\"'")
				return
			}

			var tunnelCreation = tunnel.NewTunnel()
			currentVersion := tunnelCreation.GetCurrentVersion()
			newVersion := tunnelCreation.GetNewVersion()

			if newVersion == "" {
				logger.Error("Failed to fetch the new Marijan version. Please try again later.")
				return
			}

			switch semver.Compare(currentVersion, newVersion) {
			case -1:
				// Upgrade here
				logger.Normal(fmt.Sprintf("Your current version is %s, and new version available is %s", currentVersion, newVersion))
				logger.Normal("Installing new version...")

				tunnelCreation.SetNewVersion(newVersion)

				err = tunnelCreation.Download()
				if err != nil {
					logger.Error(err.Error())
				}

				err = tunnelCreation.Extract()
				if err != nil {
					logger.Error(err.Error())
				}
				return
			case 1:
				logger.Success("You are running the latest version " + currentVersion)
				return
			case 0:
				logger.Success("You are running the latest version " + currentVersion)
				return
			default:
				logger.Success("Unknown version checker status, makesure your Marijan version is valid semver")
				return
			}
		},
	}

	return runCmd
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
