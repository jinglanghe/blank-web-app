package dao

import (
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
)

func SysVersionGet() (version model.SysVersion, err error) {
	db := GetDB().Model(&model.SysVersion{})

	result := db.First(&version)
	if result.Error != nil {
		logging.Error(result.Error).Msg("get systemversion failed")
		err = result.Error
		return
	}

	return
}
