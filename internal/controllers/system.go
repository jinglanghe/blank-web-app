package controllers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/model"
)

func registerSystemSetting(rg *gin.RouterGroup) {
	ctrl := &systemSettingController{}

	g := rg.Group("/system-setting")
	g.GET("/system-version", ctrl.GetSysVersion)
}

type systemSettingController struct {
	BaseController
}

// GetSysVersion
// @BasePath /api/v1
// @Summary 获取系统版本8888
// @Schemes
// @Description 获取系统版本
// @Tags 系统设置
// @Accept json
// @Produce json
// @Success 200 {object} model.SysVersion
// @Router /system-setting/system-version [get]
func (m *systemSettingController) GetSysVersion(c *gin.Context) {
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
