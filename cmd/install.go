package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/kubaj/doper/repository"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Docker CLI application",
	Run: func(cmd *cobra.Command, args []string) {
		err := repository.InstallPackage(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires the name of the package to install")
		}
		return nil
	},
}
