package falconClientGolang

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type DataResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type FalconClientGolang struct {
	falconClientId         string
	falconClientSecret     string
	falconUrlAccessToken   string
	falconUrlResourceOwner string
	falconAccountsUrl      string
}

func New(falconClientId, falconClientSecret, falconUrlAccessToken, falconUrlResourceOwner, falconAccountsUrl string) FalconClientGolang {
	config := FalconClientGolang{falconClientId, falconClientSecret, falconUrlAccessToken, falconUrlResourceOwner, falconAccountsUrl}
	go refreshToken(config)
	return config
}

var COOKIE_NAME string = "sdslabs"
var TOKEN string = ""

func makeRequest(url, token string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error setting up a request: %s", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making a request: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed read response: %s", err.Error())
	}
	return string(contents), nil
}

func getTokenHandler(config FalconClientGolang) {
	payload := strings.NewReader("client_id=" + config.falconClientId + "&client_secret=" + config.falconClientSecret + "&grant_type=client_credentials")
	req, _ := http.NewRequest("POST", config.falconUrlAccessToken, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	response := &DataResponse{}
	json.Unmarshal([]byte(string(body)), &response)

	TOKEN = response.AccessToken
	// return string(response.AccessToken)
}

func refreshToken(config FalconClientGolang) {
	for true {
		getTokenHandler(config)
		time.Sleep(time.Second * 3600)
	}
}

func getToken() string {
	if TOKEN == "" {
		time.Sleep(time.Second * 1)
		return TOKEN
	}
	return TOKEN
}

func GetUserById(id string, config FalconClientGolang) (string, error) {
	token := getToken()
	user_data, err := makeRequest(config.falconUrlResourceOwner+"id/"+id, token)
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	return user_data, nil
}

func GetUserByUsername(username string, config FalconClientGolang) (string, error) {
	token := getToken()
	user_data, err := makeRequest(config.falconUrlResourceOwner+"username/"+username, token)
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	return user_data, nil
}

func GetUserByEmail(email string, config FalconClientGolang) (string, error) {
	token := getToken()
	user_data, err := makeRequest(config.falconUrlResourceOwner+"email/"+email, token)
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	return user_data, nil
}

func GetLoggedInUser(config FalconClientGolang, hash string) (string, error) {
	// hash := cookies[COOKIE_NAME]
	// hash, _ := r.Cookie(COOKIE_NAME)
	// fmt.Fprint(w, cookie)
	// fmt.Println(hash)
	// var hash string = strings.Split(cookie, "=")[1].(string)
	// fmt.Println(hash)
	if hash == "" {
		return "", fmt.Errorf("cookie not found")
	}
	token := getToken()
	user_data, err := makeRequest(config.falconUrlResourceOwner+`/logged_in_user/`+hash, token)
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	return user_data, nil
}

func Login(config FalconClientGolang, w http.ResponseWriter, r *http.Request) (string, error) {
	user_data, err := GetLoggedInUser(config, "")
	if err != nil {
		return "", fmt.Errorf("failed to login with given credentials: %s", err.Error())
	}

	if user_data == "" {
		http.Redirect(w, r, config.falconAccountsUrl+`/login?redirect=//`, http.StatusTemporaryRedirect)
	}
	return user_data, nil
}
