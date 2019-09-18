package commons

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/utils"
)

// Deployer redeploys applications on different hosts
func DeployRPC(app map[string]interface{}, hostURL, service string) {
	if app["rebuild"].(bool) {
		utils.LogInfo("Re-deploying %s instance in %s\n", strings.Title(service), hostURL)
		reqBody, err := json.Marshal(app)
		if err != nil {
			utils.LogError(err)
			return
		}

		req, err := http.NewRequest("POST", "http://"+hostURL, bytes.NewBuffer(reqBody))
		if err != nil {
			utils.LogError(err)
			return
		}

		req.Header.Set("dominus-secret", configs.SWSConfig["secret"].(string))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			utils.LogError(err)
			return
		}

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		utils.LogDebug("instance has been deployed: %s", bodyString)
	}
}
