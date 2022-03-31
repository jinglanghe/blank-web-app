package controllers

import (
	"fmt"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dto"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/service/aaa"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/utils"
	"github.com/gin-gonic/gin"
)

func registerResourceQuota(rg *gin.RouterGroup) {
	ctrl := &resourceQuotaController{}

	g := rg.Group("/resource-quotas")
	g.POST("", ctrl.create)
	g.GET("", ctrl.list)
	g.GET("/:id", ctrl.get)
	g.DELETE("/:id", ctrl.delete)
}

var (
	resourceQuotaDao = &dao.ResourceQuota{}
)

type resourceQuotaController struct {
	BaseController
}

func (r *resourceQuotaController) create(c *gin.Context) {
	var createDto dto.ResourceQuotaCreate

	if !r.BindAndValidate(c, &createDto) {
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
	u, err := aaa.UserCurrent(token)
	if err != nil {
		utils.ErrUserCurrent.Message = err.Error()
		fail(c, utils.ErrUserCurrent)
		return
	}

	_type, err := nodeDeviceDao.Get(createDto.Type, createDto.Key)
	if err != nil {
		failWithDBError(c, err)
		return
	}
	if _type.Type == "" {
		fail(c, utils.ErrNodeDeviceTypeNotExist)
		return
	}

	quota := model.ResourceQuota{
		Type:        createDto.Type,
		Key:         createDto.Key,
		Arch:        _type.Arch,
		Model:       _type.Model,
		ComputeType: _type.ComputeType,
		Series:      _type.Series,
		CpuNum:      createDto.CpuNum,
		Mem:         createDto.Mem,
		CreatorID:   u.ID,
		CreatorName: u.Username,
		Num:         1,
		AvlNum:      _type.AvlNum,
	}

	exist, err = resourceQuotaDao.Exist(&quota)
	if err != nil {
		failWithDBError(c, err)
		return
	}
	if exist {
		fail(c, utils.ErrResourceQuotaAlreadyExist)
		return
	}

	if err := resourceQuotaDao.Create(&quota); err != nil {
		failWithDBError(c, err)
		return
	}

	ok(c)
}

func (r *resourceQuotaController) list(c *gin.Context) {
	var listDto dto.ResourceQuotaList
	if !r.BindAndValidate(c, &listDto) {
		return
	}

	total, err := resourceQuotaDao.Count(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	quotas, err := resourceQuotaDao.List(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	respWithData(c, map[string]interface{}{"items": quotas, "total": total})
}

func (r *resourceQuotaController) get(c *gin.Context) {
	var getDto dto.ResourceQuotaGet
	if !r.BindAndValidate(c, &getDto) {
		return
	}

	quota, err := resourceQuotaDao.Get(getDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	_type, err := nodeDeviceDao.Get(quota.Type, quota.Key)
	if err != nil {
		failWithDBError(c, err)
		return
	}
	quota.AvlNum = _type.AvlNum

	respWithData(c, quota)
}

func (r *resourceQuotaController) delete(c *gin.Context) {
	var deleteDto dto.ResourceQuotaDelete
	if !r.BindAndValidate(c, &deleteDto) {
		return
	}

	quota, err := resourceQuotaDao.Get(deleteDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	if err := resourceQuotaDao.Delete(quota); err != nil {
		failWithDBError(c, err)
		return
	}

	ok(c)
}
