package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/configs"
	"net/http"
)

type resourceConfController struct{}

func registerResourceConf(rg *gin.RouterGroup) {
	ctrl := &resourceConfController{}

	g := rg.Group("/resources-conf")
	g.GET("/metadata", ctrl.MetaData)
}

func (r *resourceConfController) MetaData(c *gin.Context) {
	types, err := nodeDeviceDao.List()
	if err != nil {
		failWithDBError(c, err)
		return
	}

	roles, err := nodeDao.Roles()
	if err != nil {
		failWithDBError(c, err)
		return
	}
	var rolesTuple []configs.ValueTuple
	for _, role := range roles {
		rolesTuple = append(rolesTuple, configs.ValueTuple{
			Value: role,
			Label: role,
		})
	}

	data := configs.QueryConditionMeta{
		Roles:      rolesTuple,
		Status:     configs.QueryConditionMetaData.Status,
		Source:     configs.QueryConditionMetaData.Source,
		Types:      types,
		AlertTypes: configs.AlertTypes(),
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})
}
