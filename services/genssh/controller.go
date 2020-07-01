// +build !windows

package genssh

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.GenSSH

// setWinsize uses low-level system call to resize the PTY device "which is just a FD in unix systems".
// See -- https://github.com/gliderlabs/ssh/blob/master/_examples/ssh-pty/pty.go
func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

// sessionHandler manages the ssh session.
func sessionHandler(s ssh.Session) {
	ptyReq, winCh, isPty := s.Pty()
	if !isPty {
		fmt.Fprintln(s, "PTY not requested")
		s.Exit(1)
		return
	}

	var cmd *exec.Cmd

	if _, err := docker.InspectContainerState(s.User()); err != nil {
		utils.LogError("GenSSH-Controller-1", err)
		utils.LogInfo("GenSSH-Controller-2", "Application %s's container not present in the current node", s.User())
		utils.LogInfo("GenSSH-Controller-3", "Attempting to a create a SSH bridge connection with the desired node")

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

// publicKeyHandler handles the public key authentication
func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	return false
}

// passwordHandler handles the password authentication
func passwordHandler(ctx ssh.Context, password string) bool {
	eventLog := "SSH login attempt `%s` on application container %s deployed at %s from IP %s"
	count, err := mongo.CountInstances(types.M{
		mongo.NameKey:         ctx.User(),
		mongo.PasswordKey:     password,
		mongo.InstanceTypeKey: mongo.AppInstance,
	})
	if err != nil {
		utils.LogInfo("GenSSH-Controller-4", "SSH login attempt failed due to unavailability of mongoDB service on host %s from IP %s", ctx.LocalAddr(), ctx.RemoteAddr())
		utils.LogError("GenSSH-Controller-5", err)
		return false
	}
	if count == 1 {
		utils.LogInfo("GenSSH-Controller-6", eventLog, "successful", ctx.User(), ctx.LocalAddr(), ctx.RemoteAddr())
		return true
	}
	utils.LogInfo("GenSSH-Controller-7", eventLog, "failed", ctx.User(), ctx.LocalAddr(), ctx.RemoteAddr())
	return false
}

// NewService returns a new instance of SSH microservice
func NewService() *ssh.Server {
	hostSigners, err := getHostSigners(configs.ServiceConfig.GenSSH.HostSigners)
	if err != nil {
		utils.Log("GenSSH-Controller-8", "There was a problem deploying GenSSH SSH service", utils.ErrorTAG)
		utils.Log("GenSSH-Controller-9", "Make sure the paths of Private Keys is correct in `config.toml`", utils.ErrorTAG)
		utils.LogError("GenSSH-Controller-10", err)
		os.Exit(1)
	}
	if !utils.IsValidPort(configs.ServiceConfig.GenSSH.Port) {
		msg := fmt.Sprintf("Port %d is invalid or already in use", configs.ServiceConfig.GenSSH.Port)
		utils.Log("GenSSH-Controller-11", msg, utils.ErrorTAG)
		os.Exit(1)
	}
	return &ssh.Server{
		Addr:             fmt.Sprintf(":%d", configs.ServiceConfig.GenSSH.Port),
		HostSigners:      hostSigners,
		Handler:          sessionHandler,
		PasswordHandler:  passwordHandler,
		PublicKeyHandler: publicKeyHandler,
	}
}
