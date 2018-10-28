package phpapp

import (
	"github.com/sdslabs/SDS/utils"
)

// readAndWriteConfig creates new config file for the given app
func (json phpAppConfig) ReadAndWriteConfig() utils.Error {
	// containerID, ok := os.LookupEnv("PHP_CONTAINER_ID")
	// if !ok {
	// 	return utils.Error{
	// 		Code: 500,
	// 		Err:  errors.New("PHP_CONTAINER_ID not found in the environment"),
	// 	}
	// }

	err := utils.ReadAndWriteConfig(json.Name, "php", "e10d2d60e777")
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
