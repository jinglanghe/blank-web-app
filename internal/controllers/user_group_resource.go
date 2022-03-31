package controllers

import (
	"fmt"
	"github.com/apulis/bmod/aistudio-aom/internal/dao"
	"github.com/apulis/bmod/aistudio-aom/internal/dto"
	"github.com/apulis/bmod/aistudio-aom/internal/model"
	"github.com/apulis/bmod/aistudio-aom/internal/service"
	"github.com/apulis/bmod/aistudio-aom/internal/service/aaa"
	"github.com/apulis/bmod/aistudio-aom/internal/utils"
	"github.com/gin-gonic/gin"
)

var (
	userGroupResourceDao = &dao.UserGroupResource{}
)

type userGroupResourceController struct {
	BaseController
}

func registerUserGroupResource(rg *gin.RouterGroup) {
	ctrl := &userGroupResourceController{}

	g := rg.Group("/user-group-resources")
	g.POST("", ctrl.create)
	g.PATCH("/:id", ctrl.update)
	g.GET("", ctrl.list)
	g.GET("/:id", ctrl.get)
	g.DELETE("/:id", ctrl.delete)
	g.GET("/:id/preview", ctrl.preview)
}

func (u *userGroupResourceController) create(c *gin.Context) {
	var createDto dto.UserGroupResourceCreate

	if !u.BindAndValidate(c, &createDto) {
		return
	}

	v, exist := c.Get("JWT_TOKEN")
	if !exist {
		err := fmt.Errorf("get jwt token failed")
		utils.ErrorMissingJwtToken.Message = err.Error()
		fail(c, utils.ErrorMissingJwtToken)
		return
	}
	token := v.(string)
	user, err := aaa.UserCurrent(token)
	if err != nil {
		utils.ErrUserCurrent.Message = err.Error()
		fail(c, utils.ErrUserCurrent)
		return
	}

	org, err := aaa.OrgDetail(token, createDto.OrgID)
	if err != nil {
		utils.ErrOrgDetail.Message = err.Error()
		fail(c, utils.ErrOrgDetail)
		return
	}
	if org.ID == 0 {
		fail(c, utils.ErrOrgDetail)
		return
	}

	userGroup, err := aaa.UserGroupDetail(token, createDto.UserGroupID)
	if err != nil {
		utils.ErrUserGroupDetail.Message = err.Error()
		fail(c, utils.ErrUserGroupDetail)
		return
	}
	if userGroup.ID == 0 {
		fail(c, utils.ErrUserGroupDetail)
		return
	}
	if userGroup.OrganizationID != org.ID {
		fail(c, utils.ErrUserGroupInvalid)
		return
	}

	if err := u.minCheck(createDto.OrgID, createDto.UserGroupID, *createDto.Min); err != nil {
		fail(c, err)
		return
	}

	resource := model.UserGroupResource{
		OrgID:         org.ID,
		OrgName:       org.Account,
		UserGroupID:   userGroup.ID,
		UserGroupName: userGroup.Account,
		Min:           *createDto.Min,
		Max:           *createDto.Max,
		Duration:      createDto.Duration,
		Source:        createDto.Source,
		CreatorID:     user.ID,
		CreatorName:   user.Username,
	}

	if err := userGroupResourceDao.Create(&resource); err != nil {
		failWithDBError(c, err)
		return
	}
	service.RefreshQuota()
	ok(c)
}

func (u *userGroupResourceController) update(c *gin.Context) {
	var updateDto dto.UserGroupResourceUpdate
	if !u.BindAndValidate(c, &updateDto) {
		return
	}

	resource, err := userGroupResourceDao.Get(updateDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	if err := u.minCheck(resource.OrgID, resource.UserGroupID, *updateDto.Min); err != nil {
		fail(c, err)
		return
	}

	resourceList, err := u.resource2ResourceList(resource)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	codeMsg := u.updateUsed(resourceList)
	if codeMsg != nil {
		fail(c, codeMsg)
		return
	}
	if *updateDto.Min < resourceList.Used {
		utils.ErrUserGroupResourceMinInvalid.Message = "min can not less then used"
		fail(c, utils.ErrUserGroupResourceMinInvalid)
		return
	}

	resource.Min = *updateDto.Min
	resource.Max = *updateDto.Max
	if err := userGroupResourceDao.Update(resource); err != nil {
		failWithDBError(c, err)
		return
	}

	service.RefreshQuota()
	ok(c)
}

func (u *userGroupResourceController) list(c *gin.Context) {
	var listDto dto.UserGroupResourceList
	if !u.BindAndValidate(c, &listDto) {
		return
	}

	total, err := userGroupResourceDao.Count(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	resources, err := userGroupResourceDao.List(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	nodeDeviceTypeOrgs, err := nodeDeviceDao.OrgDevises()
	if err != nil {
		failWithDBError(c, err)
		return
	}

	orgAvlMemes, err := nodeDao.OrgAvlMemes()
	if err != nil {
		failWithDBError(c, err)
		return
	}

	resourceLists := userGroupResourceDao.MatchResources(resources, nodeDeviceTypeOrgs, orgAvlMemes)
	for k, list := range resourceLists {
		codeMsg := u.updateUsed(&list)
		if codeMsg != nil {
			fail(c, codeMsg)
			return
		}
		resourceLists[k] = list
	}

	respWithData(c, map[string]interface{}{"items": resourceLists, "total": total})
}

func (u *userGroupResourceController) get(c *gin.Context) {
	var getDto dto.UserGroupResourceGet
	if !u.BindAndValidate(c, &getDto) {
		return
	}

	resource, err := userGroupResourceDao.Get(getDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	resourceList, err := u.resource2ResourceList(resource)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	codeMsg := u.updateUsed(resourceList)
	if codeMsg != nil {
		fail(c, codeMsg)
		return
	}

	respWithData(c, resourceList)
}

func (u *userGroupResourceController) delete(c *gin.Context) {
	var deleteDto dto.UserGroupResourceDelete
	if !u.BindAndValidate(c, &deleteDto) {
		return
	}

	resource, err := userGroupResourceDao.Get(deleteDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	if err := userGroupResourceDao.Delete(resource); err != nil {
		failWithDBError(c, err)
		return
	}

	ok(c)
}

func (u *userGroupResourceController) preview(c *gin.Context) {
	var previewDto dto.UserGroupResourcePreview
	if !u.BindAndValidate(c, &previewDto) {
		return
	}

	resource, err := userGroupResourceDao.Get(previewDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	if err := u.minCheck(resource.OrgID, resource.UserGroupID, *previewDto.Min); err != nil {
		fail(c, err)
		return
	}

	nodeDeviceTypes, err := nodeDeviceDao.OrgDevs(resource.OrgID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	orgAvlMem, err := nodeDao.OrgAvlMem(resource.OrgID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	resource.Min = *previewDto.Min
	resource.Max = *previewDto.Max
	resourceList := userGroupResourceDao.MatchResource(resource, nodeDeviceTypes, orgAvlMem)

	codeMsg := u.updateUsed(resourceList)
	if codeMsg != nil {
		fail(c, codeMsg)
		return
	}
	if *previewDto.Min < resourceList.Used {
		utils.ErrUserGroupResourceMinInvalid.Message = "min can not less then used"
		fail(c, utils.ErrUserGroupResourceMinInvalid)
		return
	}

	respWithData(c, resourceList)
}

func (u *userGroupResourceController) resource2ResourceList(resource *model.UserGroupResource) (*model.UserGroupResourceList, error) {
	nodeDeviceTypes, err := nodeDeviceDao.OrgDevs(resource.OrgID)
	if err != nil {
		return nil, err
	}

	orgAvlMem, err := nodeDao.OrgAvlMem(resource.OrgID)
	if err != nil {
		return nil, err
	}

	list := userGroupResourceDao.MatchResource(resource, nodeDeviceTypes, orgAvlMem)

	return list, nil
}

func (u *userGroupResourceController) updateUsed(list *model.UserGroupResourceList) *utils.CodeMessage {
	allocatedReses, err := service.NamespaceAllocatedReses(list.UserGroupName)
	if err != nil {
		utils.ErrNamespaceUsedRes.Message = err.Error()
		return utils.ErrNamespaceUsedRes
	}

	var used int64

	for _, allocatedRes := range allocatedReses {
		if allocatedRes.IsMem() {
			memUsed := utils.DividedUp(allocatedRes.Value(), list.Mem)
			used = utils.Max(used, memUsed)
			continue
		}

		for key, quota := range list.Quotas {
			if allocatedRes.EqualQuota(&quota) {
				list.Quotas[key].Used = allocatedRes.Value()
				resUsed := utils.DividedUp(allocatedRes.Value(), quota.Num)
				used = utils.Max(used, resUsed)
				break
			}
		}
	}

	list.Used = used
	return nil
}

func (u *userGroupResourceController) minCheck(orgId, userGroupId, min int64) *utils.CodeMessage {
	errCodeMessage := utils.ErrUserGroupResourceExceed

	sumMin, err := userGroupResourceDao.SumMin(orgId, userGroupId)
	if err != nil {
		errCodeMessage.Message = err.Error()
		return errCodeMessage
	}
	if sumMin+min > 100 {
		return errCodeMessage
	}
	return nil
}
