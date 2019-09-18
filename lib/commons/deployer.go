package commons

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/utils"
)

// Deployer redeploys applications on different hosts
func Deployer(app map[string]interface{}, hostURL, service string) {
	if app["rebuild"].(bool) {
		reqBody, err := json.Marshal(app)
		if err != nil {
			utils.LogError(err)
		}

		req, err := http.NewRequest("POST", hostURL, bytes.NewBuffer(reqBody))
		if err != nil {
			utils.LogError(err)
		}

		req.Header.Set("dominus-secret", configs.SWSConfig["secret"].(string))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			utils.LogError(err)
		}

		defer resp.Body.Close()
	}
}
