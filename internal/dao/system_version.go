package dao

import (
	"github.com/apulis/bmod/aistudio-aom/internal/model"
	"github.com/apulis/sdk/go-utils/logging"
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
