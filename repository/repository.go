package repository

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kubaj/doper/env"
)

var (
	coreRepository = "https://github.com/kubaj/doper-packages.git"
)

func ResolvePackage(pkg string) string {
	s := strings.Split(pkg, "/")
	if len(s) == 1 {
		return ""
	}

	return pkg
}

func InstallPackage(pkg string) error {

	err := InitPackageDir(env.PackagePath)
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "doper")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	err = exec.Command("git", "clone", coreRepository, dir).Run()
	if err != nil {
		return err
	}

	config := filepath.Join(dir, pkg, "doper.yaml")
	_, err = os.Stat(config)
	if os.IsNotExist(err) {
		fmt.Println("Package does not exist")
		return err
	}

	err = exec.Command("cp", config, filepath.Join(env.PackagePath, pkg+".yaml")).Run()
	return err
}

func InitPackageDir(packageDir string) error {
	stat, err := os.Stat(packageDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err = os.MkdirAll(packageDir, 0755)
		return err
	}

	if !stat.IsDir() {
		return errors.New(packageDir + " is not a directory")
	}

	return nil
}
