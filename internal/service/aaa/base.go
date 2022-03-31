package aaa

import (
	"fmt"
	config "github.com/apulis/bmod/aistudio-aom/configs"
)

func getAaaDomain() string {
	aaa := fmt.Sprintf("%s:%d", config.Config.IAM.Host, config.Config.IAM.Port)
	return aaa
}
