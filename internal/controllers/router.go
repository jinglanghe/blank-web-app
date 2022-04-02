package controllers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	config "gitlab.apulis.com.cn/hjl/blank-web-app/configs"
	docs "gitlab.apulis.com.cn/hjl/blank-web-app/docs"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/jwt"
)

var (
	jwtAuth *jwt.AuthN
)

func RegisterRoutes(e *gin.Engine) {
	e.MaxMultipartMemory = 8 << 20 // 8 MiB
	docs.SwaggerInfo.BasePath = config.Config.APIPrefix
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	jwtConfig := config.Config.JWTConfig

	jwtAuth = jwt.NewJwtAuthN(
		jwt.SigningAlgorithm("HS256"),
		jwt.SecretKey([]byte("jwt secret key")),
		jwt.PublicKey(jwtConfig.PublicKeyFile),
	)
	v1 := e.Group(config.Config.APIPrefix)

	registerSystemSetting(v1)

	v1.Use(jwtAuth.Middleware())

}
