package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kubaj/doper/config"

	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	packageDir = "packages"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Docker CLI application",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		b, err := ioutil.ReadFile(filepath.Join(home, ".doper", packageDir, args[0]+".yaml"))
		if err != nil {
			fmt.Errorf("package %s not installed", args[0])
			os.Exit(1)
		}
		p := config.Package{}
		err = yaml.Unmarshal(b, &p)
		if err != nil {
			fmt.Errorf("failed to read config")
			os.Exit(1)
		}

		p.Run()

	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires the name of the package to run")
		}
		return nil
	},
}
