package main

import (
	"github.com/devetek/d-panel-cli/internal/api"
	"github.com/devetek/d-panel-cli/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type AuthCmd struct {
	cmd       *cobra.Command
	zapLogger *zap.Logger

	email    string
	password string
}

func NewAuthCmd(logger *zap.Logger) *AuthCmd {
	return &AuthCmd{
		zapLogger: logger,
		cmd: &cobra.Command{
			Use:   "auth",
			Short: "Manage dPanel session",
		},
	}
}

func (u *AuthCmd) Connect() *cobra.Command {
	u.cmd.AddCommand(
		u.login(),
	)

	return u.cmd
}

func (u *AuthCmd) login() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "login",
		Short: "Authorize to access dPanel",
		Run: func(cmd *cobra.Command, args []string) {
			if u.email == "" || u.password == "" {
				u.zapLogger.Error("Email and password are required")
				return
			}

			// init dPanel client
			client := api.NewClient()

			// check if session exist
			err := client.CheckSessionExist()
			if err != nil {
				_, err := client.Login(u.email, u.password)
				if err != nil {
					logger.Error("Login error: " + err.Error())
					return
				}

				// double check profile
				_, err = client.GetProfile()
				if err != nil {
					logger.Error("Error get user profile: " + err.Error())
					return
				}
			}

			logger.Success("Success login to dPanel!")
		},
	}

	runCmd.PersistentFlags().StringVarP(&u.email, "email", "e", "", "Youre registered email address in dPanel")
	runCmd.PersistentFlags().StringVarP(&u.password, "password", "p", "", "Youre registered password in dPanel")

	return runCmd
}
