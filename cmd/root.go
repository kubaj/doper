package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Find home directory.
	// home, err := homedir.Dir()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// viper.AddConfigPath(home)
	// viper.SetConfigName(".doper")

	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Println("Can't read config:", err)
	// 	os.Exit(1)
	// }
}

var rootCmd = &cobra.Command{
	Use:   "doper",
	Short: "Doper is package manager for Docker CLI applications",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
