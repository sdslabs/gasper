package app

import (
	"io/ioutil"
	"strings"

	"github.com/sdslabs/SDS/utils"
)

// readAndWriteStaticConf creates new config file for static app
func readAndWriteStaticConf(name string) utils.Error {
	content, err := ioutil.ReadFile("")
	if err != nil {
		return utils.Error{
			Code: 500,
			Err:  err,
		}
	}

	conf := strings.Replace(string(content), "template", name, -1)

	newContent := []byte(conf)
	err = ioutil.WriteFile(""+name+".static.sdslabs.co.conf", newContent, 0644)
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
