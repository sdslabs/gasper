package gin

import (
	"errors"
	"fmt"
)

var immutableFields = []string{"name", "_id", "instanceType", "containerID", "execID", "hostIP", "httpPort", "language"}

func validateUpdatePayload(data map[string]interface{}) error {
	res := ""
	for _, field := range immutableFields {
		if data[field] != nil {
			res += fmt.Sprintf("Field `%s` is immutable; ", field)
		}
	}
	if res != "" {
		return errors.New(res)
	}
	return nil
}
