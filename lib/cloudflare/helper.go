package cloudflare

import (
	"errors"
	"fmt"
)

func formatErrorResponse(apiErrors []errorResponse) error {
	res := ""
	for _, value := range apiErrors {
		res += value.Message + ";"
	}
	return fmt.Errorf("Cloudflare Errors: %s", res)
}

// getZoneID returns the ID of the 1st zone
func getZoneID() (string, error) {
	res, err := GetZones()
	if err != nil {
		return "", err
	}
	if len(res.Result) == 0 {
		return "", errors.New("No active zones available at the moment")
	}
	return res.Result[0].ID, err
}

// getRecordID returns the ID of the desired record
func getRecordID(name, instanceType string) (string, error) {
	res, err := FetchRecords(map[string]interface{}{
		"name": fmt.Sprintf("%s.%s.%s", name, instanceType, domain),
	})
	if err != nil {
		return "", err
	}
	if len(res.Result) == 0 {
		return "", fmt.Errorf("Domain name for application %s of type %s doesn't exist", name, instanceType)
	}
	return res.Result[0].ID, err
}
