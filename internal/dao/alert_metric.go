package dao

import (
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AlertMetric struct{}

func (r *AlertMetric) Inc(_uuid uuid.UUID, _type string) error {
	alertMetric := model.AlertMetric{
		UUID:  _uuid,
		Type:  _type,
		Count: 1,
	}
	db := GetDB().Model(&model.AlertMetric{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"count": gorm.Expr("alert_metrics.count + 1")}),
	}).Create(&alertMetric)
	return db.Error
}

func (r *AlertMetric) ListTypes() ([]model.AlertMetricType, error) {
	var types []model.AlertMetricType
	db := GetDB().Model(&model.AlertMetric{}).
		Select("type, sum(count) as count").
		Group("type").
		Find(&types)
	return types, db.Error
}
