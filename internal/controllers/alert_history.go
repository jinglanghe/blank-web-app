package controllers

import (
	"github.com/apulis/bmod/aistudio-aom/internal/dao"
	_ "github.com/apulis/bmod/aistudio-aom/internal/dao"
	"github.com/apulis/bmod/aistudio-aom/internal/dto"
	"github.com/apulis/go-business/pkg/jwt"

	//"github.com/apulis/bmod/aistudio-aom/internal/session"
	"github.com/gin-gonic/gin"
)

var (
	serviceAlertsDao  = &dao.ServiceAlerts{}
	resourceAlertsDao = &dao.ResourceAlerts{}
)

func registerAlertHistory(rg *gin.RouterGroup) {
	ctrl := &alertHistoryController{}

	serviceHistoryGroup := rg.Group("service-history")
	serviceHistoryGroup.GET("", ctrl.listService)
	serviceHistoryGroup.DELETE("", ctrl.deleteServiceAlerts)
	serviceHistoryGroup.PUT("/:uuid", ctrl.serviceAlertChangeStatus)

	resourceHistoryGroup := rg.Group("resource-history")
	resourceHistoryGroup.GET("", ctrl.listResource)
	resourceHistoryGroup.DELETE("", ctrl.deleteResourceAlerts)
	resourceHistoryGroup.PUT("/:uuid", ctrl.resourceAlertChangeStatus)
}

type alertHistoryController struct {
	BaseController
}

func (a *alertHistoryController) listResource(c *gin.Context) {
	var listDto dto.ResourceAlertListDto
	if !a.BindAndValidate(c, &listDto) {
		return
	}

	orgId, _ := jwt.OrgId(c)
	listDto.OrgId = orgId

	alerts, err := resourceAlertsDao.List(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}
	total, err := resourceAlertsDao.Count(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	respWithData(c, map[string]interface{}{
		"total": total,
		"items": alerts,
	})
}

func (a *alertHistoryController) deleteResourceAlerts(c *gin.Context) {
	var deleteDto dto.ResourceAlertsDeleteDto
	if !a.BindAndValidate(c, &deleteDto) {
		return
	}

	if err := resourceAlertsDao.Delete(deleteDto.UUIDs); err != nil {
		failWithDBError(c, err)
		return
	}
	ok(c)
}

func (a *alertHistoryController) resourceAlertChangeStatus(c *gin.Context) {
	var statusDto dto.ResourceAlertStatusDto
	if !a.BindAndValidate(c, &statusDto) {
		return
	}

	alert, err := resourceAlertsDao.Get(statusDto.UUID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	alert.Status = *statusDto.Status
	if err := resourceAlertsDao.Update(alert); err != nil {
		failWithDBError(c, err)
		return
	}
	ok(c)
}

func (a *alertHistoryController) listService(c *gin.Context) {
	var listDto dto.ServiceAlertListDto
	if !a.BindAndValidate(c, &listDto) {
		return
	}

	orgId, _ := jwt.OrgId(c)
	listDto.OrgId = orgId

	alerts, err := serviceAlertsDao.List(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}
	total, err := serviceAlertsDao.Count(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	respWithData(c, map[string]interface{}{
		"total": total,
		"items": alerts,
	})
}

func (a *alertHistoryController) deleteServiceAlerts(c *gin.Context) {
	var deleteDto dto.ServiceAlertsDeleteDto
	if !a.BindAndValidate(c, &deleteDto) {
		return
	}

	if err := serviceAlertsDao.Delete(deleteDto.UUIDs); err != nil {
		failWithDBError(c, err)
		return
	}
	ok(c)
}

func (a *alertHistoryController) serviceAlertChangeStatus(c *gin.Context) {
	var statusDto dto.ServiceAlertStatusDto
	if !a.BindAndValidate(c, &statusDto) {
		return
	}

	alert, err := serviceAlertsDao.Get(statusDto.UUID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	alert.Status = *statusDto.Status
	if err := serviceAlertsDao.Update(alert); err != nil {
		failWithDBError(c, err)
		return
	}
	ok(c)
}
