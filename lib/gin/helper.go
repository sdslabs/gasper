package gin

import (
	"errors"
	"fmt"

	"github.com/sdslabs/gasper/types"
)

var immutableFields = []string{"name", "_id", "instanceType", "containerID", "execID", "hostIP", "httpPort", "language", "cloudflareID", "domainURL"}

func validateUpdatePayload(data types.M) error {
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
