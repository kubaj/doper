package env

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	ConfigPath  string
	PackagePath string
)

func init() {
	c, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ConfigPath = filepath.Join(c, ".doper")
	PackagePath = filepath.Join(ConfigPath, "packages")
}
