package ssh

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

const (
	// DefaultServiceName is the name of the SSH microservice
	DefaultServiceName = types.SSH

	// ProxyServiceName is the name of the proxy service of SSH
	ProxyServiceName = types.SSHProxy
)

// newHandler returns a handler function which manages the ssh session.
func newHandler(service string) func(s ssh.Session) {
	var proxy bool
	if service == ProxyServiceName {
		proxy = true
	} else if service == DefaultServiceName {
		proxy = false
	} else {
		return nil
	}
	return func(s ssh.Session) {
		ptyReq, winCh, isPty := s.Pty()
		if !isPty {
			fmt.Fprintln(s, "PTY not requested")
			s.Exit(1)
			return
		}

		var cmd *exec.Cmd

		if proxy {
			instanceURL, err := redis.FetchAppNode(s.User())
			if err != nil {
				fmt.Fprintln(s, fmt.Sprintf("Application %s is not deployed at the moment", s.User()))
				s.Exit(1)
				return
			}
			instanceURL = strings.Split(instanceURL, ":")[0]
			port, err := redis.GetSSHPort(instanceURL)
			if err != nil {
				fmt.Fprintln(s, "Sorry, we are experiencing some technical difficulties at the moment")
				s.Exit(1)
				return
			}
			if port == redis.ErrEmptySet {
				fmt.Fprintln(s, fmt.Sprintf("Instance %s doesn't have the SSH service deployed", instanceURL))
				s.Exit(1)
				return
			}
			cmd = exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-p", port, fmt.Sprintf("%s@%s", s.User(), instanceURL))
		} else {
			cmd = exec.Command("docker", "exec", "-it", s.User(), "/bin/sh")
		}
		cmd.Env = append(cmd.Env, s.Environ()...)
		termEnv := fmt.Sprintf("TERM=%s", ptyReq.Term)
		cmd.Env = append(cmd.Env, termEnv)

		ptmx, err := pty.Start(cmd)
		if err != nil {
			fmt.Fprintf(s, "ERROR: %s", err.Error())
			s.Exit(1)
			return
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
	return false
}

// passwordHandler handles the password authentication
func passwordHandler(ctx ssh.Context, password string) bool {
	eventLog := "SSH login attempt `%s` on application container %s deployed at %s from IP %s"
	count, err := mongo.CountInstances(types.M{
		"name":         ctx.User(),
		"password":     password,
		"instanceType": mongo.AppInstance,
	})
	if err != nil {
		utils.LogInfo("SSH login attempt failed due to unavailability of mongoDB service on host %s from IP %s", ctx.LocalAddr(), ctx.RemoteAddr())
		utils.LogError(err)
		return false
	}
	if count == 1 {
		utils.LogInfo(eventLog, "successful", ctx.User(), ctx.LocalAddr(), ctx.RemoteAddr())
		return true
	}
	utils.LogInfo(eventLog, "failed", ctx.User(), ctx.LocalAddr(), ctx.RemoteAddr())
	return false
}

func serviceBuilder(service string, filepaths []string, port int) (*ssh.Server, error) {
	hostSigners, err := getHostSigners(filepaths)
	if err != nil {
		return nil, err
	}
	if !utils.IsValidPort(port) {
		msg := fmt.Sprintf("Port %d is invalid or already in use.\n", port)
		utils.Log(msg, utils.ErrorTAG)
		return nil, errors.New(msg)
	}
	return &ssh.Server{
		Addr:             fmt.Sprintf(":%d", port),
		HostSigners:      hostSigners,
		Handler:          newHandler(service),
		PasswordHandler:  passwordHandler,
		PublicKeyHandler: publicKeyHandler,
	}, nil
}

func handleError(err error) {
	utils.Log("There was a problem deploying SSH service", utils.ErrorTAG)
	utils.Log("Make sure the paths of Private Keys is correct in `config.json`", utils.ErrorTAG)
	utils.LogError(err)
	panic(err)
}

// NewDefaultService returns a new instance of SSH microservice
func NewDefaultService() *ssh.Server {
	server, err := serviceBuilder(DefaultServiceName, configs.ServiceConfig.SSH.HostSigners, configs.ServiceConfig.SSH.Port)
	if err != nil {
		handleError(err)
	}
	return server
}

// NewProxyService returns a new proxy instance of SSH microservice
func NewProxyService() *ssh.Server {
	server, err := serviceBuilder(ProxyServiceName, configs.ServiceConfig.SSHProxy.HostSigners, configs.ServiceConfig.SSHProxy.Port)
	if err != nil {
		handleError(err)
	}
	return server
}
