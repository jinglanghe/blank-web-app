package controllers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/model"
)

func registerMetric(rg *gin.RouterGroup) {
	ctrl := &metricController{}

	g := rg.Group("/metrics")
	g.GET("/system-version", ctrl.GetSysVersion)
}

type metricController struct {
	BaseController
}

func (m *metricController) GetSysVersion(c *gin.Context) {
	sysVersion := model.SysVersion{}

	// 从数据库
	sysVersion, err := dao.SysVersionGet()
	if err == nil {
		respWithData(c, sysVersion)
		return
	}

	// 从配置文件
	fileContent, err := ioutil.ReadFile("../internal/metadata/sys_version.json")
	if err == nil && len(fileContent) > 0 {
		err = json.Unmarshal(fileContent, &sysVersion)
		if err == nil {
			respWithData(c, sysVersion)
			return
		}
	}

	sysVersion.SysVersion = "1.0.0"
	sysVersion.Installer = "admin"
	sysVersion.InstallTime = 1623910517706
	sysVersion.Description = "enjoy yourself"
	respWithData(c, sysVersion)
}
