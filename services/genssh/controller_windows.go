package genssh

import (
	"github.com/gliderlabs/ssh"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.GenSSH

// GenSSH doesn't work on windows, hence this function returns a nil instance
func NewService() *ssh.Server {
	return nil
}
