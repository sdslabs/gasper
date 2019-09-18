package commons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sdslabs/SWS/lib/configs"
)

// Deployer redeploys applications on different hosts
func Deployer(app map[string]interface{}, url, service string) {
	hostURL := url + configs.ServiceConfig[service].(map[string]interface{})["port"].(string)
	reqBody, err := json.Marshal(app)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reqBody)
	resp, err := http.Post(hostURL,
		"application/json",
		bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
}
