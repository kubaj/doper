package config

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

type Config struct {
	Packages map[string]string
}

type Package struct {
	Image        string
	Volumes      []string
	Ports        []string
	PreserveUser bool `yaml:"preserveUser"`
	Env          map[string]string
	Entrypoint   string
	Workdir      string
}

func (p *Package) Run() error {
	if err := p.ExpandVolumes(); err != nil {
		return err
	}

	args := []string{"run", "-it"}

	if p.PreserveUser {
		uid, gid := p.GetIDs()
		args = append(args, "-u", uid+":"+gid)
	}

	for _, v := range p.Volumes {
		args = append(args, "-v", v)
	}

	for _, p := range p.Ports {
		args = append(args, "-p", p)
	}

	for key, val := range p.Env {
		args = append(args, "-e", key+"="+val)
	}

	if p.Entrypoint != "" {
		args = append(args, "--entrypoint", p.Entrypoint)
	}

	if p.Workdir != "" {
		args = append(args, "-w", p.Workdir)
	}

	args = append(args, p.Image)

	cmd := exec.Command("docker", args...)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}

func (p *Package) GetIDs() (string, string) {
	uid := strconv.Itoa(os.Getuid())
	gid := strconv.Itoa(os.Getgid())
	return uid, gid
}

func (p *Package) ExpandVolumes() error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for i := range p.Volumes {
		p.Volumes[i] = strings.Replace(p.Volumes[i], "$(pwd)", dir, -1)
	}

	return nil
}
