package controllers

import (
	"encoding/json"
	"fmt"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/service"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func GetSysVersion(c *gin.Context) {
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

func DownloadCerts(c *gin.Context) {
	data, err := service.DownloadCerts()
	if err != nil {
		fail(c, err)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=cert.zip")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	c.Writer.Write(data)
}
