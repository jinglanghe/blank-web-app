package aaa

import (
	"fmt"
	config "gitlab.apulis.com.cn/hjl/blank-web-app-2/configs"
)

func getAaaDomain() string {
	aaa := fmt.Sprintf("%s:%d", config.Config.IAM.Host, config.Config.IAM.Port)
	return aaa
}
