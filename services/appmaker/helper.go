package appmaker

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

var path, _ = os.Getwd()

// storageCleanup removes the application's local storage directory
func storageCleanup(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		utils.LogError("AppMaker-Helper-1", err)
	}
	return err
}

// containerCleanup removes the application's container
func containerCleanup(appName string) error {
	err := docker.DeleteContainer(appName)
	if err != nil {
		utils.LogError("AppMaker-Helper-2", err)
	}
	return err
}

// diskCleanup cleans the specified application's container and local storage
func diskCleanup(appName string) {
	appDir := filepath.Join(path, fmt.Sprintf("storage/%s", appName))
	storeCleanupChan := make(chan error)
	go func() {
		storeCleanupChan <- storageCleanup(appDir)
	}()
	containerCleanup(appName)
	<-storeCleanupChan
}

// stateCleanup removes the application's data from MongoDB and Redis
func stateCleanup(appName string) {
	_, err := mongo.DeleteInstance(types.M{
		mongo.NameKey:         appName,
		mongo.InstanceTypeKey: mongo.AppInstance,
	})
	if err != nil {
		utils.LogError("AppMaker-Helper-3", err)
	}
	if err := redis.RemoveApp(appName); err != nil {
		utils.LogError("AppMaker-Helper-4", err)
	}
}

func generateSSHKeys() ([]byte, []byte) {
	// filename := "key"
	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	// Extract public component.
	pub := key.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.
	pvtPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)
	return pvtPEM, pubPEM
}
