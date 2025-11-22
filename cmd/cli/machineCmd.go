package main

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/devetek/d-panel-cli/internal/api"
	"github.com/devetek/d-panel-cli/internal/helper"
	"github.com/devetek/d-panel-cli/internal/logger"
	"github.com/devetek/d-panel-cli/internal/tunnel"
	"github.com/devetek/d-panel/pkg/dmachine"
	"github.com/devetek/d-panel/pkg/drouter"
	"github.com/devetek/d-panel/pkg/dsecret"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type MachineCmd struct {
	cmd       *cobra.Command
	zapLogger *zap.Logger

	sshIP    string
	sshPort  string
	httpPort string
	// use behind tunnel if this machine is behind NAT without public IP
	behindTunnel bool
	// used when your machine want to expose dPanel agent behind proxy
	// this option used by dPanel to access your machine, example input:
	// - http://my-machine-01.devetek.app -> for insecure connection / HTTP
	// - https://my-machine-01.devetek.app -> for Secure connection / HTTPS
	domain string
}

func NewMachineCmd(logger *zap.Logger) *MachineCmd {
	return &MachineCmd{
		zapLogger: logger,
		cmd: &cobra.Command{
			Use:   "machine",
			Short: "Manage dPanel machine",
		},
	}
}

func (m *MachineCmd) Connect() *cobra.Command {
	m.cmd.AddCommand(
		m.create(),
	)

	return m.cmd
}

func (m *MachineCmd) create() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "create",
		Short: "Add this machine to dPanel",
		Long:  `Add this machine to dPanel and manage easily.`,
		Run: func(cmd *cobra.Command, args []string) {
			// check if user has sudo access in golang
			if !helper.IsSudo() {
				logger.Error("You must run this command as sudo, currenty dpanel-agent required to running under root")
				return
			}

			// init dPanel client
			client := api.NewClient()

			// check if session exist
			err := client.CheckSessionExist()
			if err != nil {
				logger.Error("Please login to your dPanel account, use command 'dnocs auth login --email=\"email@email.com\" --password=\"password\"'")
				return
			}

			// get list secret ssh
			secretSSH, err := client.GetListSecretSSH()
			if err != nil {
				logger.Error("Error get list secret ssh: " + err.Error())
				return
			}

			var mySSHKey dsecret.Response
			if secretSSH.Data.Pagination.TotalItem == 0 {
				// create new SSH key
				newSSHKey, err := client.CreateSecretSSH()
				if err != nil {
					logger.Error("Error create secret ssh: " + err.Error())
					return
				}

				// assign from new SSH key
				mySSHKey = newSSHKey.Data
			} else {
				// get first secret ssh from existing
				mySSHKey = secretSSH.Data.Secrets[0]

				// get detail secret ssh
				detailSSHKey, err := client.GetSecretSSHByID(fmt.Sprintf("%d", mySSHKey.ID))
				if err != nil {
					logger.Error("Error get detail secret ssh: " + err.Error())
					return
				}

				// assign from detail SSH key
				mySSHKey = detailSSHKey.Data
			}

			if !helper.IsSSHAuthorized(mySSHKey.Data.Data()["public"]) {
				// append ssh key to authorized_keys file
				err = helper.AppendAuthorizedKey(mySSHKey.Data.Data()["public"])
				if err != nil {
					logger.Error("Error append ssh key to authorized_keys file: " + err.Error())
					return
				}
			}

			if !m.behindTunnel {
				// make sure sshIP is not empty
				if m.sshIP == "" {
					// get my public IP automatically
					m.sshIP, err = helper.GetMyIP()
					if err != nil {
						logger.Error("Error get my public IP " + err.Error())
						return
					}
				}

				if m.httpPort == "" {
					// get available port
					availablePort, err := helper.FindAvailablePort()
					if err != nil {
						logger.Error("Error get available port + " + err.Error())
						return
					}

					m.httpPort = fmt.Sprintf("%d", availablePort)
				}
			}

			currentUser, err := user.Current()
			if err != nil {
				logger.Error("Error getting current user: " + err.Error())
				return
			}

			// integrate with tunnel
			if m.behindTunnel {
				var currentTunnel = tunnel.NewTunnel()

				// check tunnel configs
				var tunnelConfig = currentTunnel.GetConfig()
				if len(tunnelConfig) == 0 {
					logger.Error("This machine is not connected to dPanel tunnel")
					return
				}

				var tunnelHTTPPort string
				var originHTTPPort string
				for _, tunnel := range tunnelConfig {
					if strings.Contains(tunnel.ID, "ssh-") {
						m.sshIP = tunnel.TunnelHost
						m.sshPort = tunnel.ListenerPort
					}

					if strings.Contains(tunnel.ID, "http-") {
						tunnelHTTPPort = tunnel.ListenerPort
						originHTTPPort = tunnel.ServicePort
					}
				}

				// set payload
				var payload = drouter.PayloadRouter{
					AdvanceMode: false,
					Type:        "proxy_pass",
					Name:        fmt.Sprintf("http-%s-to-%s", tunnelHTTPPort, originHTTPPort),
					Domain:      fmt.Sprintf("http-%s-to-%s 1", tunnelHTTPPort, originHTTPPort),
					MachineID:   11,
					Upstream:    fmt.Sprintf("localhost:%s", tunnelHTTPPort),
				}

				router, err := client.CreateRouter(payload)
				if err != nil {
					logger.Error("Failed to create HTTP server for this machine, with error " + err.Error())
					logger.Error("Login to dPanel, open https://cloud-beta.terpusat.com/router, and delete existing domain")
					return
				}

				// set domain for this machine
				m.httpPort = originHTTPPort
				m.domain = router.Data.Domain
			}

			// register new server
			newServer := dmachine.Payload{
				Provider: "other",
				SecretID: fmt.Sprintf("%d", mySSHKey.ID),
				Address:  m.sshIP,
				SSHPort:  m.sshPort,
				HTTPPort: m.httpPort,
				Domain:   m.domain,
				SSHUser:  currentUser.Username,
			}

			// check if server already registered
			if client.IsRegistered() {
				logger.Error("Server already registered with your account")
				return
			}

			// register new server
			server, err := client.RegisterServer(newServer)
			if err != nil {
				logger.Error("Error register server " + err.Error())
				return
			}

			// setup server
			_, err = client.SetupServer(int(server.Data.ID))
			if err != nil {
				logger.Error("Error setup server " + err.Error())
				return
			}

			logger.Success("Success register server, visit " + api.FrontendURL + "/v2/resources/servers to check the progress!")
		},
	}

	runCmd.PersistentFlags().StringVarP(&m.sshIP, "ssh-ip", "i", "", "SSH IP of your machine")
	runCmd.PersistentFlags().StringVarP(&m.sshPort, "ssh-port", "s", "22", "SSH port of your machine")
	runCmd.PersistentFlags().StringVarP(&m.httpPort, "http-port", "p", "9000", "HTTP port of your machine")
	runCmd.PersistentFlags().StringVarP(&m.domain, "http-domain", "d", "", "HTTP domain of agent (optional)")
	runCmd.PersistentFlags().BoolVarP(&m.behindTunnel, "behind-tunnel", "t", false, "Read tunnel config and auto create domain")

	return runCmd
}
