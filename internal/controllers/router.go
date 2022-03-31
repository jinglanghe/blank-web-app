package controllers

import (
	config "github.com/apulis/bmod/aistudio-aom/configs"
	"github.com/apulis/go-business/pkg/jwt"
	"github.com/gin-gonic/gin"
)

var (
	jwtAuth *jwt.AuthN
)

func RegisterRoutes(e *gin.Engine) {
	e.MaxMultipartMemory = 8 << 20 // 8 MiB

	jwtConfig := config.Config.JWTConfig

	jwtAuth = jwt.NewJwtAuthN(
		jwt.SigningAlgorithm("HS256"),
		jwt.SecretKey([]byte("jwt secret key")),
		jwt.PublicKey(jwtConfig.PublicKeyFile),
	)
	v1 := e.Group("/api/v1")


	registerMetric(v1)

	v1.Use(jwtAuth.Middleware())

	v1.GET("/system-version", GetSysVersion)
	v1.GET("/system/certs", DownloadCerts)
	v1.POST("/modelarts-setting", SetModelArts)
	v1.GET("/modelarts-setting", ModelArtsList)

	rg := v1.Group("/alert")

	registerAlertHistory(rg)
	registerNode(v1)
	registerResourceConf(v1)
	registerOrgResource(v1)
	registerUserGroupResource(v1)
	registerResourceQuota(v1)
	registerPrometheus(v1)
}
