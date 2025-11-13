package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/devetek/d-panel-cli/internal/api"
	"github.com/devetek/d-panel-cli/internal/helper"
	"github.com/devetek/d-panel/pkg/dmachine"
	"github.com/devetek/d-panel/pkg/dsecret"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// User input variables
var email string
var password string

var sshIP string
var sshPort string

var httpPort string

func runCmd() *cobra.Command {
	// init zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logger.Sync()

	var runCmd = &cobra.Command{
		Use:   "register",
		Short: "Register this machine to dPanel",
		Long:  `Register this machine to dPanel, this will create a new tunnel client on your machine.`,
		Run: func(cmd *cobra.Command, args []string) {
			// check if user has sudo access in golang
			// if !helper.IsSudo() {
			// 	logger.Error("You must run this command as sudo")
			// 	return
			// }

			// email and password is required
			if email == "" || password == "" {
				logger.Error("Email and password are required")
				return
			}

			// init dPanel client
			client := api.NewClient()

			// check if session exist
			err = client.CheckSessionExist()
			if err != nil {
				err = client.Login(email, password)
				if err != nil {
					logger.Error("Error login to dPanel", zap.Error(err))
					return
				}
			}

			// get user profile
			_, err = client.GetProfile()
			if err != nil {
				logger.Error("Error get user profile", zap.Error(err))
				return
			}

			// get list secret ssh
			secretSSH, err := client.GetListSecretSSH()
			if err != nil {
				logger.Error("Error get list secret ssh", zap.Error(err))
				return
			}

			var mySSHKey dsecret.Response
			if secretSSH.Data.Pagination.TotalItem == 0 {
				// create new SSH key
				newSSHKey, err := client.CreateSecretSSH()
				if err != nil {
					logger.Error("Error create secret ssh", zap.Error(err))
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
					logger.Error("Error get detail secret ssh", zap.Error(err))
					return
				}

				// assign from detail SSH key
				mySSHKey = detailSSHKey.Data
			}

			// append ssh key to authorized_keys file
			err = helper.AppendAuthorizedKey(mySSHKey.Data.Data()["public"])
			if err != nil {
				logger.Error("Error append ssh key to authorized_keys file", zap.Error(err))
				return
			}

			// make sure sshIP is not empty
			if sshIP == "" {
				// get my public IP automatically
				sshIP, err = helper.GetMyIP()
				if err != nil {
					logger.Error("Error get my public IP", zap.Error(err))
					return
				}
			}

			if httpPort == "" {
				// get available port
				availablePort, err := helper.FindAvailablePort()
				if err != nil {
					logger.Error("Error get available port", zap.Error(err))
					return
				}

				httpPort = fmt.Sprintf("%d", availablePort)
			}

			currentUser, err := user.Current()
			if err != nil {
				fmt.Printf("Error getting current user: %v\n", err)
				return
			}

			// register new server
			newServer := dmachine.Payload{
				Provider: "other",
				SecretID: fmt.Sprintf("%d", mySSHKey.ID),
				Address:  sshIP,
				SSHPort:  sshPort,
				HTTPPort: httpPort,
				Domain:   "",
				SSHUser:  currentUser.Username,
			}

			// register new server
			server, err := client.RegisterServer(newServer)
			if err != nil {
				logger.Error("Error register server", zap.Error(err))
				return
			}

			// setup server
			setup, err := client.SetupServer(int(server.Data.ID))
			if err != nil {
				logger.Error("Error setup server", zap.Error(err))
				return
			}

			logger.Info("Success setup server", zap.Any("setup", setup))
		},
	}

	runCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "Youre registered email address in dPanel")
	runCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Youre registered password in dPanel")
	runCmd.PersistentFlags().StringVarP(&sshIP, "ssh-ip", "i", "", "SSH IP of your machine")
	runCmd.PersistentFlags().StringVarP(&sshPort, "ssh-port", "s", "22", "SSH port of your machine")
	runCmd.PersistentFlags().StringVarP(&httpPort, "http-port", "t", "9000", "HTTP port of your machine")

	runCmd.MarkPersistentFlagRequired("email")
	runCmd.MarkPersistentFlagRequired("password")

	return runCmd
}
