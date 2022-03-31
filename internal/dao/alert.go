package dao

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"

	"github.com/apulis/bmod/aistudio-aom/internal/dto"
	"github.com/apulis/bmod/aistudio-aom/internal/model"
	"github.com/apulis/sdk/go-utils/logging"
	"gorm.io/datatypes"
)

type ResourceAlerts struct{}

func (r *ResourceAlerts) Create(record *model.ResourceAlert) error {
	return GetDB().Save(record).Error
}

func (r *ResourceAlerts) Get(uuid string) (*model.ResourceAlert, error) {
	var alert model.ResourceAlert
	db := GetDB().Where("uuid = ?", uuid).First(&alert)
	return &alert, db.Error
}

func (r *ResourceAlerts) Count(cond *dto.ResourceAlertListDto) (int64, error) {
	db := GetDB().Model(&model.ResourceAlert{}).Where("org_id = ?", cond.OrgId).
		Where("created_at BETWEEN ? AND ?", cond.StartTime, cond.EndTime)
	if cond.Level > 0 {
		db = db.Where("level = ?", cond.Level)
	}
	if len(cond.Name) > 0 {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", cond.Name))
	}
	var count int64
	db = db.Count(&count)

	return count, db.Error
}

func (r *ResourceAlerts) List(cond *dto.ResourceAlertListDto) ([]model.ResourceAlert, error) {
	db := GetDB().Model(&model.ResourceAlert{}).Where("org_id = ?", cond.OrgId).
		Where("created_at BETWEEN ? AND ?", cond.StartTime, cond.EndTime)
	if cond.Level > 0 {
		db = db.Where("level = ?", cond.Level)
	}
	if len(cond.Name) > 0 {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", cond.Name))
	}

	var alerts []model.ResourceAlert
	db = db.Offset((cond.BaseListDto.PageNum - 1) * cond.BaseListDto.PageSize).
		Limit(cond.BaseListDto.PageSize).Order("created_at DESC").
		Find(&alerts)

	for k, alert := range alerts {
		alerts[k].Receivers = convertReceivers(alert.Receivers)
	}

	return alerts, db.Error
}

func (r *ResourceAlerts) AppendReceiver(uuid string, result map[string]bool) error {
	tx := GetDB().Begin()

	data, _ := json.Marshal(result)
	err := tx.Model(&model.ResourceAlert{}).Where("uuid = ?", uuid).
		Update("receivers", gorm.Expr("coalesce(receivers::jsonb, '{}'::jsonb) || (?)::jsonb", string(data))).Error
	if err != nil {
		logging.Error(err).Msg("update alert receivers error")
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *ResourceAlerts) Delete(uuids []string) error {
	return GetDB().Where("uuid IN ?", uuids).Delete(model.ResourceAlert{}).Error
}

func (r *ResourceAlerts) Update(record *model.ResourceAlert) error {
	return GetDB().Select("*").Updates(record).Error
}

type ServiceAlerts struct{}

func (s *ServiceAlerts) Get(uuid string) (*model.ServiceAlert, error) {
	var alert model.ServiceAlert
	db := GetDB().Where("uuid = ?", uuid).First(&alert)
	return &alert, db.Error
}

func (s *ServiceAlerts) Create(record *model.ServiceAlert) error {
	return GetDB().Save(record).Error
}

func (s *ServiceAlerts) Count(cond *dto.ServiceAlertListDto) (int64, error) {
	db := GetDB().Model(&model.ServiceAlert{}).Where("org_id = ?", cond.OrgId).
		Where("created_at BETWEEN ? AND ?", cond.StartTime, cond.EndTime)
	if len(cond.Name) > 0 {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", cond.Name))
	}
	var count int64
	db = db.Count(&count)

	return count, db.Error
}

func (s *ServiceAlerts) List(cond *dto.ServiceAlertListDto) ([]model.ServiceAlert, error) {
	db := GetDB().Model(&model.ServiceAlert{}).Where("org_id = ?", cond.OrgId).
		Where("created_at BETWEEN ? AND ?", cond.StartTime, cond.EndTime)
	if len(cond.Name) > 0 {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", cond.Name))
	}

	var alerts []model.ServiceAlert
	db = db.Offset((cond.BaseListDto.PageNum - 1) * cond.BaseListDto.PageSize).
		Limit(cond.BaseListDto.PageSize).Order("created_at DESC").
		Find(&alerts)
	for k, alert := range alerts {
		alerts[k].Receivers = convertReceivers(alert.Receivers)
	}

	return alerts, db.Error
}

func (s *ServiceAlerts) AppendReceiver(uuid string, result map[string]bool) error {
	tx := GetDB().Begin()

	data, _ := json.Marshal(result)
	err := tx.Model(&model.ServiceAlert{}).Where("uuid = ?", uuid).
		Update("receivers", gorm.Expr("coalesce(receivers::jsonb, '{}'::jsonb) || (?)::jsonb", string(data))).Error
	if err != nil {
		logging.Error(err).Msg("update alert receivers error")
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *ServiceAlerts) Delete(uuids []string) error {
	return GetDB().Where("uuid IN ?", uuids).Delete(model.ServiceAlert{}).Error
}

func (s *ServiceAlerts) Update(record *model.ServiceAlert) error {
	return GetDB().Select("*").Updates(record).Error
}

type receivers struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
}

func convertReceivers(rawData datatypes.JSON) datatypes.JSON {
	m := make(map[string]bool)
	_ = json.Unmarshal(rawData, &m)

	r := receivers{Success: []string{}, Failed: []string{}}
	for k, v := range m {
		if v {
			r.Success = append(r.Success, k)
		} else {
			r.Failed = append(r.Failed, k)
		}
	}
	d, _ := json.Marshal(r)
	return d
}
