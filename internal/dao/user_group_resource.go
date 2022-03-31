package dao

import (
	"database/sql"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dto"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/service/aaa"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/utils"
	"github.com/apulis/sdk/go-utils/logging"
	"gorm.io/gorm"
)

type UserGroupResource struct{}

func CreateSystemAdminGroupResource() {
	systemAdminGroup := aaa.SystemAdminGroup()

	userGroupResource := model.UserGroupResource{
		UserGroupID:   systemAdminGroup.ID,
		UserGroupName: systemAdminGroup.Account,
	}

	if err := createUserGroupResource(&userGroupResource); err != nil {
		logging.Error(err).Msg("create system admin group resources failed ")
	}
}

func CreateDefaultUserGroupResource() {
	defaultUserGroup := aaa.DefaultUserGroup()

	userGroupResource := model.UserGroupResource{
		OrgID:         defaultUserGroup.OrganizationID,
		OrgName:       defaultUserGroup.Organization.Account,
		UserGroupID:   defaultUserGroup.ID,
		UserGroupName: defaultUserGroup.Account,
		Min:           100,
		Max:           100,
	}
	if err := createUserGroupResource(&userGroupResource); err != nil {
		logging.Error(err).Msg("create default user group resources failed ")
	}
}

func createUserGroupResource(userGroupResource *model.UserGroupResource) error {
	var count int64
	GetDB().Model(&model.UserGroupResource{}).Where("user_group_id = ?", userGroupResource.UserGroupID).Count(&count)
	if count > 0 {
		return nil
	}
	db := GetDB().Create(&userGroupResource)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (u *UserGroupResource) Create(resource *model.UserGroupResource) error {
	return GetDB().Save(resource).Error
}

func (u *UserGroupResource) Get(id int64) (*model.UserGroupResource, error) {
	var resource model.UserGroupResource
	db := GetDB().Where("id = ?", id).First(&resource)
	return &resource, db.Error
}

func (u *UserGroupResource) Count(cond *dto.UserGroupResourceList) (int64, error) {
	var count int64
	db := GetDB().Model(&model.UserGroupResource{})
	db = u.QueryFilter(db, cond)
	db = db.Count(&count)

	return count, db.Error
}

func (u *UserGroupResource) List(cond *dto.UserGroupResourceList) ([]model.UserGroupResource, error) {
	var resources []model.UserGroupResource
	db := GetDB().Model(&model.UserGroupResource{})
	db = u.QueryFilter(db, cond)
	db = sort(db, cond.Sort)
	db = db.Offset((cond.BaseListDto.PageNum - 1) * cond.BaseListDto.PageSize).
		Limit(cond.BaseListDto.PageSize).Order("created_at DESC").
		Find(&resources)

	return resources, db.Error
}

func (u *UserGroupResource) QueryFilter(db *gorm.DB, cond *dto.UserGroupResourceList) *gorm.DB {
	if cond.OrgId != 0 {
		db = db.Where("org_id = ?", cond.OrgId)
	}
	if cond.UserGroupId != 0 {
		db = db.Where("user_group_id = ?", cond.UserGroupId)
	}
	return db
}

func (u *UserGroupResource) Update(resource *model.UserGroupResource) error {
	return GetDB().Model(resource).Select("min", "max").Updates(*resource).Error
}

func (u *UserGroupResource) Delete(resource *model.UserGroupResource) error {
	return GetDB().Delete(&resource).Error
}

func (u *UserGroupResource) SumMin(orgId, userGroupId int64) (int64, error) {
	var sum sql.NullInt64
	db := GetDB().Model(&model.UserGroupResource{}).
		Select("sum(min)").
		Where("org_id = ?", orgId).
		Where("user_group_id != ?", userGroupId).
		Scan(&sum)
	return sum.Int64, db.Error
}

func (u *UserGroupResource) ListALl() ([]model.UserGroupResource, error) {
	var resources []model.UserGroupResource
	db := GetDB().Model(&model.UserGroupResource{})
	db = db.Find(&resources)

	return resources, db.Error
}

func (u *UserGroupResource) MatchResources(resources []model.UserGroupResource, nodeDeviceTypeOrgs []model.NodeDeviceTypeOrg, orgAvlMems []model.OrgAvlMem) (list []model.UserGroupResourceList) {
	for _, resource := range resources {
		orgAvlMem := u.findOrgAvlMem(orgAvlMems, resource.OrgID)
		nodeDeviceTypes := u.findNodeDeviceTypes(nodeDeviceTypeOrgs, resource.OrgID)
		l := u.MatchResource(&resource, nodeDeviceTypes, orgAvlMem)
		list = append(list, *l)
	}
	return list
}

func (u *UserGroupResource) findOrgAvlMem(orgAvlMems []model.OrgAvlMem, orgId int64) int64 {
	for _, orgAvlMem := range orgAvlMems {
		if orgAvlMem.OrgId == orgId {
			return orgAvlMem.AvlMem
		}
	}
	return 0
}

func (u *UserGroupResource) findNodeDeviceTypes(nodeDeviceTypeOrgs []model.NodeDeviceTypeOrg, orgId int64) []model.NodeDeviceType {
	var nodeDeviceTypes []model.NodeDeviceType
	for _, nodeDeviceTypeOrg := range nodeDeviceTypeOrgs {
		if nodeDeviceTypeOrg.OrgId == orgId {
			nodeDeviceType := model.NodeDeviceType{
				Key:         nodeDeviceTypeOrg.Key,
				Type:        nodeDeviceTypeOrg.Type,
				Arch:        nodeDeviceTypeOrg.Arch,
				Model:       nodeDeviceTypeOrg.Model,
				ComputeType: nodeDeviceTypeOrg.ComputeType,
				Series:      nodeDeviceTypeOrg.Series,
				Num:         nodeDeviceTypeOrg.Num,
				AvlNum:      nodeDeviceTypeOrg.AvlNum,
			}
			nodeDeviceTypes = append(nodeDeviceTypes, nodeDeviceType)
		}
	}
	return nodeDeviceTypes
}

func (u *UserGroupResource) MatchResource(resource *model.UserGroupResource, nodeDeviceTypes []model.NodeDeviceType, orgAvlMem int64) *model.UserGroupResourceList {
	list := model.UserGroupResourceList{
		UserGroupResource: *resource,
		Mem:               orgAvlMem,
		MemMin:            utils.RoundDown(orgAvlMem, resource.Min),
		MemMax:            utils.RoundUp(orgAvlMem, resource.Max),
		Quotas:            nil,
	}

	for _, nodeDeviceType := range nodeDeviceTypes {
		quota := model.Quota{
			Type:        nodeDeviceType.Type,
			Key:         nodeDeviceType.Key,
			Arch:        nodeDeviceType.Arch,
			Model:       nodeDeviceType.Model,
			ComputeType: nodeDeviceType.ComputeType,
			Series:      nodeDeviceType.Series,
			Num:         nodeDeviceType.AvlNum,
			Min:         utils.RoundDown(nodeDeviceType.AvlNum, resource.Min),
			Max:         utils.RoundUp(nodeDeviceType.AvlNum, resource.Max),
		}
		list.Quotas = append(list.Quotas, quota)
	}
	return &list
}
