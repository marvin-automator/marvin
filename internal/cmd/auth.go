package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/marvin-automator/marvin/internal"
	"github.com/marvin-automator/marvin/internal/auth"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

var setPassword = &cobra.Command{
	Use:   "set_password",
	Short: "Set the password to log in to marvin using the web interface.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter a password: ")

		fd := int(syscall.Stdin)
		pw, err := terminal.ReadPassword(fd)
		internal.ErrorAndExit(err)

		err = auth.SetPassword(string(pw))
		internal.ErrorAndExit(err)
		color.Green("Password saved")
	},
}
