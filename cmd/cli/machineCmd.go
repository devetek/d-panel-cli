package main

import (
	"fmt"
	"os/user"

	"github.com/devetek/d-panel-cli/internal/api"
	"github.com/devetek/d-panel-cli/internal/helper"
	"github.com/devetek/d-panel-cli/internal/logger"
	"github.com/devetek/d-panel/pkg/dmachine"
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
				logger.Error("You must run this command as sudo")
				return
			}

			// init dPanel client
			client := api.NewClient()

			// check if session exist
			err := client.CheckSessionExist()
			if err != nil {
				logger.Error("Error check session exist: " + err.Error())
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

			currentUser, err := user.Current()
			if err != nil {
				logger.Error("Error getting current user: " + err.Error())
				return
			}

			// register new server
			newServer := dmachine.Payload{
				Provider: "other",
				SecretID: fmt.Sprintf("%d", mySSHKey.ID),
				Address:  m.sshIP,
				SSHPort:  m.sshPort,
				HTTPPort: m.httpPort,
				Domain:   "",
				SSHUser:  currentUser.Username,
			}

			// TODO: Create file ~/.devetek/machine.json to prevent multiple registration

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
	runCmd.PersistentFlags().StringVarP(&m.httpPort, "http-port", "t", "9000", "HTTP port of your machine")

	return runCmd
}
