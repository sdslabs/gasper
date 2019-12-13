package iwa

import (
	"io/ioutil"

	"github.com/gliderlabs/ssh"
	"github.com/sdslabs/gasper/configs"
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

	if configs.ServiceConfig.Iwa.UsingPassphrase {
		signer, err = gossh.ParsePrivateKeyWithPassphrase(
			key,
			[]byte(configs.ServiceConfig.Iwa.Passphrase),
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
