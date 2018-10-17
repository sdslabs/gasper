package staticapp

import (
	"github.com/sdslabs/SDS/utils"
)

// readAndWriteConfig creates new config file for the given app
func (json staticAppConfig) ReadAndWriteConfig() utils.Error {
	// containerID, ok := os.LookupEnv("STATIC_CONTAINER_ID")
	// if !ok {
	// 	return utils.Error{
	// 		Code: 500,
	// 		Err:  errors.New("STATIC_CONTAINER_ID not found in the environment"),
	// 	}
	// }

	err := utils.ReadAndWriteConfig(json.Name, "static", "6f36625894e5")
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
