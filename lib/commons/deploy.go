package commons

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/utils"
)

// DeployRPC re-deploys applications on different hosts
func DeployRPC(app map[string]interface{}, hostURL, service string) {
	utils.LogInfo("Re-deploying application %s with type %s to %s\n", app["name"], strings.Title(service), hostURL)

	app["rebuild"] = true
	app["context"] = app["context"].(primitive.D).Map()
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
	utils.LogDebug("Application %s with type %s has been succesfully re-deployed to %s", app["name"], strings.Title(service), hostURL)
	utils.LogDebug("Response: %s", bodyString)
}
