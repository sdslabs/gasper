package ssh

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/utils"
)

// newHandler returns a handler function which manages the ssh session.
func newHandler(service string) func(s ssh.Session) {
	var proxy bool
	if service == "ssh_proxy" {
		proxy = true
	} else {
		proxy = false
	}
	return func(s ssh.Session) {
		ptyReq, winCh, isPty := s.Pty()
		if !isPty {
			fmt.Fprintln(s, "PTY not requested")
			s.Exit(1)
		}

		var cmd *exec.Cmd

		if proxy {
			instanceURL, err := redis.FetchAppURL(s.User())
			if err != nil {
				fmt.Fprintln(s, "Application "+s.User()+" is not deployed at the moment")
				s.Exit(1)
				return
			}
			instanceURL = strings.Split(instanceURL, ":")[0]
			cmd = exec.Command("ssh", "-p", "2222", fmt.Sprintf("%s@%s", s.User(), instanceURL))
		} else {
			cmd = exec.Command("docker", "exec", "-it", s.User(), "/bin/bash")
		}
		cmd.Env = append(cmd.Env, s.Environ()...)
		termEnv := fmt.Sprintf("TERM=%s", ptyReq.Term)
		cmd.Env = append(cmd.Env, termEnv)

		ptmx, err := pty.Start(cmd)
		if err != nil {
			fmt.Fprintf(s, "ERROR: %s", err.Error())
			s.Exit(1)
		}
		defer ptmx.Close()

		go func() {
			for win := range winCh {
				setWinsize(ptmx, win.Width, win.Height)
			}
		}()

		go func() {
			io.Copy(ptmx, s) // STDIN
		}()
		io.Copy(s, ptmx) // STDOUT
	}
}

// publicKeyHandler handles the public key authentication
func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	return true
}

// passwordHandler handles the password authentication
func passwordHandler(ctx ssh.Context, password string) bool {
	return true
}

// BuildSSHServer creates a server for the given parameters
func BuildSSHServer(service string) (*ssh.Server, error) {
	sshConfig := utils.ServiceConfig[service].(map[string]interface{})
	filepaths := utils.ToStringSlice(sshConfig["host_signers"])
	hostSigners, err := getHostSigners(service, filepaths)
	if err != nil {
		return nil, err
	}
	return &ssh.Server{
		Addr:        sshConfig["port"].(string),
		HostSigners: hostSigners,

		Handler:          newHandler(service),
		PasswordHandler:  passwordHandler,
		PublicKeyHandler: publicKeyHandler,
	}, nil
}
