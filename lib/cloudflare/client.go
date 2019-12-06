package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sdslabs/gasper/types"
)

// GetZones returns information of all zones associated
// with the domain name
func GetZones() (*ZoneResponse, error) {
	req, _ := http.NewRequest("GET", listZonesEndpoint, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	query := req.URL.Query()
	query.Add("name", domain)
	query.Add("status", "active")
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	data := &ZoneResponse{}

	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}

	if !data.Success {
		return nil, formatErrorResponse(data.Errors)
	}
	return data, nil
}

// FetchRecords returns all the DNS records for the given zone
func FetchRecords(queryParams ...types.M) (*MultiResponse, error) {
	zoneID, err := getZoneID()
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf(fetchRecordEndpoint, zoneID), nil)
	req.Header.Add("Authorization", "Bearer "+token)

	query := req.URL.Query()
	for _, params := range queryParams {
		for key, value := range params {
			query.Add(key, value.(string))
		}
	}
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	data := &MultiResponse{}

	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}

	if !data.Success {
		return nil, formatErrorResponse(data.Errors)
	}
	return data, nil
}

// createRecord creates a new DNS record for the given zone and returns its ID
func createRecord(name, instanceType string) (*SingleResponse, error) {
	zoneID, err := getZoneID()
	if err != nil {
		return nil, err
	}

	payload := &singlePayload{
		Name:    fmt.Sprintf("%s.%s", name, instanceType),
		Type:    "A",
		Content: publicIP,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf(createRecordEndpoint, zoneID), bytes.NewBuffer(payloadBytes))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	data := &SingleResponse{}

	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}

	if !data.Success {
		// If entry already exists then update the record
		res, err := updateRecord(name, instanceType, payload)
		if err != nil {
			return nil, err
		}
		if !res.Success {
			return nil, formatErrorResponse(append(data.Errors, res.Errors...))
		}
		return res, nil
	}
	return data, nil
}

// CreateApplicationRecord creates a new DNS record for an application
func CreateApplicationRecord(name string) (*SingleResponse, error) {
	return createRecord(name, ApplicationInstance)
}

// updateRecord updates the DNS record for an application in the given zone
func updateRecord(name, instanceType string, payload *singlePayload) (*SingleResponse, error) {
	zoneID, err := getZoneID()
	if err != nil {
		return nil, err
	}

	recordID, err := getRecordID(name, instanceType)
	if err != nil {
		return nil, err
	}

	payload.Type = "A"
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("PUT", fmt.Sprintf(updateRecordEndpoint, zoneID, recordID), bytes.NewBuffer(payloadBytes))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	data := &SingleResponse{}

	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}

	if !data.Success {
		return nil, formatErrorResponse(data.Errors)
	}
	return data, nil
}

// DeleteRecord deletes the DNS record for an application in the given zone
func DeleteRecord(name, instanceType string) (*GenericResponse, error) {
	zoneID, err := getZoneID()
	if err != nil {
		return nil, err
	}

	recordID, err := getRecordID(name, instanceType)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf(deleteRecordEndpoint, zoneID, recordID), nil)
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	data := &GenericResponse{}

	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}

	if !data.Success {
		return nil, formatErrorResponse(data.Errors)
	}
	return data, nil
}
