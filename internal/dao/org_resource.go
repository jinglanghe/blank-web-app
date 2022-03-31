package dao

import (
	"github.com/apulis/bmod/aistudio-aom/internal/dto"
	"github.com/apulis/bmod/aistudio-aom/internal/model"
	"github.com/apulis/bmod/aistudio-aom/internal/service/aaa"
	"github.com/apulis/sdk/go-utils/logging"
	"gorm.io/gorm"
)

type OrgResource struct{}

func CreateDefaultOrgResource() {
	defaultOrg := aaa.DefaultOrg()
	var count int64
	GetDB().Model(&model.OrgResource{}).Where("org_id = ?", defaultOrg.ID).Count(&count)
	if count > 0 {
		return
	}

	newOrgResource := model.OrgResource{
		OrgID:   defaultOrg.ID,
		OrgName: defaultOrg.Account,
	}
	db := GetDB().Create(&newOrgResource)
	if db.Error != nil {
		logging.Error(db.Error).Msg("org resource create default org resources failed ")
	}
}

func DefaultOrgResource() model.OrgResource {
	defaultOrg := aaa.DefaultOrg()

	var o model.OrgResource
	GetDB().Model(&model.OrgResource{}).Where("org_id = ?", defaultOrg.ID).First(&o)
	return o
}

func (o *OrgResource) Save(orgResource *model.OrgResource) error {
	db := GetDB().Session(&gorm.Session{FullSaveAssociations: true}).Save(&orgResource)
	if db.Error != nil {
		logging.Error(db.Error).Msg("org resource refresh failed at save org resource")
		return db.Error
	}
	return nil
}

func (o *OrgResource) Create(orgResource *model.OrgResource) error {
	return GetDB().Create(orgResource).Error
}

func (o *OrgResource) Exist(oldNode model.Node, nodeIds []int64) bool {
	for _, nodeId := range nodeIds {
		if oldNode.ID == nodeId {
			return true
		}
	}
	return false
}

func (o *OrgResource) Count(cond *dto.OrgResourceList) (int64, error) {
	var count int64
	db := GetDB().Model(&model.OrgResource{}).Count(&count)

	return count, db.Error
}

func (o *OrgResource) List(cond *dto.OrgResourceList) ([]model.OrgResource, error) {
	var orgResources []model.OrgResource
	db := GetDB().Preload("Nodes").Model(&model.OrgResource{})
	db = sort(db, cond.Sort)
	db = db.Offset((cond.BaseListDto.PageNum - 1) * cond.BaseListDto.PageSize).
		Limit(cond.BaseListDto.PageSize).Order("created_at DESC").
		Find(&orgResources)
	return orgResources, db.Error
}

func (o *OrgResource) Get(id int64) (*model.OrgResource, error) {
	var orgResource model.OrgResource
	db := GetDB().Preload("Nodes").Where("id = ?", id).First(&orgResource)
	return &orgResource, db.Error
}

func (o *OrgResource) GetByOrgId(orgId int64) (*model.OrgResource, error) {
	var orgResource model.OrgResource
	db := GetDB().Model(&model.OrgResource{}).Preload("Nodes").Where("org_id = ?", orgId).Limit(1).Find(&orgResource)
	return &orgResource, db.Error
}

func (o *OrgResource) Delete(orgResource *model.OrgResource) error {
	return GetDB().Delete(&orgResource).Error
}
