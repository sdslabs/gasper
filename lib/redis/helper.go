package redis

import (
	"encoding/json"

	"github.com/sdslabs/gasper/types"
)

// fetchBindings returns a struct containing an instance's server and node URL
func fetchBindings(key, name string) (*types.InstanceBindings, error) {
	result, err := client.HGet(key, name).Result()
	if err != nil {
		return nil, err
	}
	instance := &types.InstanceBindings{}
	resultByte := []byte(result)
	err = json.Unmarshal(resultByte, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// fetchServer returns the server URL of an instance
func fetchServer(key, name string) (string, error) {
	instance, err := fetchBindings(key, name)
	if err != nil {
		return "", err
	}
	return instance.Server, nil
}

// fetchNode returns the node URL of an instance
func fetchNode(key, name string) (string, error) {
	instance, err := fetchBindings(key, name)
	if err != nil {
		return "", err
	}
	return instance.Node, nil
}
