package dao

import (
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dto"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ResourceQuota struct{}

func (r *ResourceQuota) Refresh(_types []model.NodeDeviceType) error {
	resourceQuotas, err := r.ListAll()
	if err != nil {
		return err
	}

	var updateResourceQuotas []model.ResourceQuota
	for _, resourceQuota := range resourceQuotas {
		for _, _type := range _types {
			if resourceQuota.EqualNodeDeviceType(&_type) &&
				resourceQuota.AvlNum != _type.AvlNum {
				resourceQuota.AvlNum = _type.AvlNum
				updateResourceQuotas = append(updateResourceQuotas, resourceQuota)
			}
		}
	}

	db := GetDB().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"avl_num"}),
	}).CreateInBatches(updateResourceQuotas, 10)
	return db.Error
}

func (r *ResourceQuota) ListAll() ([]model.ResourceQuota, error) {
	var resourceQuotas []model.ResourceQuota
	db := GetDB().Model(&model.ResourceQuota{}).Find(&resourceQuotas)
	return resourceQuotas, db.Error
}

func (r *ResourceQuota) Create(resourceQuota *model.ResourceQuota) error {
	return GetDB().Save(resourceQuota).Error
}

func (r *ResourceQuota) Get(id int64) (*model.ResourceQuota, error) {
	var resourceQuota model.ResourceQuota
	db := GetDB().Model(&model.ResourceQuota{}).Where("id = ?", id).First(&resourceQuota)
	return &resourceQuota, db.Error
}

func (r *ResourceQuota) Count(cond *dto.ResourceQuotaList) (int64, error) {
	var count int64
	db := GetDB().Model(&model.ResourceQuota{})
	db = r.QueryFilter(db, cond)
	db = db.Count(&count)

	return count, db.Error
}

func (r *ResourceQuota) List(cond *dto.ResourceQuotaList) ([]model.ResourceQuota, error) {
	var resourceQuotas []model.ResourceQuota
	db := GetDB().Model(&model.ResourceQuota{})
	db = r.QueryFilter(db, cond)
	db = sort(db, cond.Sort)
	db = db.Offset((cond.BaseListDto.PageNum - 1) * cond.BaseListDto.PageSize).
		Limit(cond.BaseListDto.PageSize).Order("created_at DESC").
		Find(&resourceQuotas)

	return resourceQuotas, db.Error
}

func (r *ResourceQuota) QueryFilter(db *gorm.DB, cond *dto.ResourceQuotaList) *gorm.DB {
	if cond.Type != "" {
		db = db.Where("type = ?", cond.Type)
	}
	if cond.CpuNum != 0 {
		db = db.Where("cpu_num >= ?", cond.CpuNum)
	}
	if cond.Mem != 0 {
		db = db.Where("mem >= ?", cond.Mem)
	}
	if cond.Arch != "" {
		db = db.Where("arch = ?", cond.Arch)
	}
	if cond.CreatorName != "" {
		cond.CreatorName = "%" + cond.CreatorName + "%"
		db = db.Where("creator_name like ?", cond.CreatorName)
	}
	if cond.AvlNum != 0 {
		db = db.Where("avl_num >= ?", cond.AvlNum)
	}
	return db
}

func (r *ResourceQuota) Delete(resourceQuota *model.ResourceQuota) error {
	return GetDB().Delete(&resourceQuota).Error
}

func (r *ResourceQuota) Exist(resourceQuota *model.ResourceQuota) (bool, error) {
	var count int64
	db := GetDB().Model(&model.ResourceQuota{}).Where(resourceQuota, "type", "key", "arch", "model", "compute_type", "series", "cpu_num", "mem", "num").Count(&count)
	if db.Error != nil {
		return false, db.Error
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
