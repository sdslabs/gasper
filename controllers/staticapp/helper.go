package staticapp

import (
	"io/ioutil"
	"strings"

	"github.com/sdslabs/SDS/utils"
)

// readAndWriteConfig creates new config file for the given app
func (json appConfig) ReadAndWriteConfig() utils.Error {
	content, err := ioutil.ReadFile("")
	if err != nil {
		return utils.Error{
			Code: 500,
			Err:  err,
		}
	}

	conf := strings.Replace(string(content), "template", json.Name, -1)

	newContent := []byte(conf)
	err = ioutil.WriteFile(""+json.Name+".static.sdslabs.co.conf", newContent, 0644)
	if err != nil {
		return utils.Error{
			Code: 500,
			Err:  err,
		}
	}

	return utils.Error{
		Code: 200,
		Err:  nil,
	}
}
