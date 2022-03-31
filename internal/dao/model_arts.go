package dao

import (
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dto"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"github.com/apulis/sdk/go-utils/logging"
	"gorm.io/gorm/clause"
)

func ModelArtsCreate(m *model.ModelArts) error {
	// 有就更新，没有就插入
	result := GetDB().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&m)

	if result.Error != nil {
		logging.Error(result.Error).Msg("create model arts failed")
		return result.Error
	}

	return nil
}

func ModelArtsList(m *dto.ModelArts) (ms []model.ModelArts, count int64, err error) {
	db := GetDB().Model(&model.ModelArts{})

	result := db.Count(&count)
	if result.Error != nil {
		logging.Error(result.Error).Msg("get model arts count failed")
		err = result.Error
		return
	}

	result = db.Find(&ms)
	if result.Error != nil {
		logging.Error(result.Error).Msg("get model arts list failed")
		err = result.Error
		return
	}

	return
}
