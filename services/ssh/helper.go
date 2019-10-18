package ssh

import (
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/mongo"
	gossh "golang.org/x/crypto/ssh"
)

// getPrivateKey returns a Signer interface for the private key
// specified from the filepath
func getPrivateKey(filepath string) (ssh.Signer, error) {

	var err error
	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var signer gossh.Signer

	if configs.ServiceConfig.SSH.UsingPassphrase {
		signer, err = gossh.ParsePrivateKeyWithPassphrase(
			key,
			[]byte(configs.ServiceConfig.SSH.Passphrase),
		)
	} else {
		signer, err = gossh.ParsePrivateKey(key)
	}

	if err != nil {
		return nil, err
	}

	return signer.(ssh.Signer), nil
}

// getHostSigners returns a slice of Signer interface for the
// specified filepaths of the private keys
func getHostSigners(filepaths []string) ([]ssh.Signer, error) {
	var signers []ssh.Signer
	for _, filepath := range filepaths {
		signer, err := getPrivateKey(filepath)
		if err != nil {
			return nil, err
		}
		signers = append(signers, signer)
	}
	return signers, nil
}

// getPublicKey returns a PublicKey interface for the key
// specified from the filepath
func getPublicKey(filepath string) (ssh.PublicKey, error) {
	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	publicKey, err := gossh.ParsePublicKey(key)
	if err != nil {
		return nil, err
	}

	return publicKey.(ssh.PublicKey), nil
}

// setWinsize uses low-level system call to resize the PTY device "which is just a FD in unix systems".
// See -- https://github.com/gliderlabs/ssh/blob/master/_examples/ssh-pty/pty.go
func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

// getAppFromContext gets fetches the app docs for the ssh user
func getAppFromContext(ctx ssh.Context) map[string]interface{} {
	app := mongo.FetchAppInfo(map[string]interface{}{
		"name": ctx.User(),
	})
	// Should return just one app
	return app[0]
}
