package ssh

import (
	"github.com/gliderlabs/ssh"
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
