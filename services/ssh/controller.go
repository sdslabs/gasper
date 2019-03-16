package ssh

import (
	"github.com/gliderlabs/ssh"
	"github.com/sdslabs/SWS/lib/utils"
)

// Handler handles the ssh session.
func Handler(s ssh.Session) {}

// PublicKeyHandler handles the public key authentication
func PublicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	return true
}

// PasswordHandler handles the password authentication
func PasswordHandler(ctx ssh.Context, password string) bool {
	return true
}

// BuildSSHServer creates a server for the given parameters
func BuildSSHServer() (*ssh.Server, error) {
	sshConfig := utils.ServiceConfig["ssh"].(map[string]interface{})
	filepaths := utils.ToStringSlice(sshConfig["host_signers"])
	hostSigners, err := getHostSigners(filepaths)
	if err != nil {
		return nil, err
	}
	return &ssh.Server{
		Addr:        sshConfig["port"].(string),
		HostSigners: hostSigners,

		Handler:          Handler,
		PasswordHandler:  PasswordHandler,
		PublicKeyHandler: PublicKeyHandler,
	}, nil
}
