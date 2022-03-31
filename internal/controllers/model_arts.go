package controllers

import (
	"encoding/base64"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dto"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/utils"
	"github.com/gin-gonic/gin"
)

func SetModelArts(c *gin.Context) {
	dtoInput := dto.ModelArts{}
	if b := bindAndValidate(c, &dtoInput); !b {
		return
	}

	encryptSK, err := utils.AesEncrypt([]byte(dtoInput.SK), utils.DefaultKey)
	if err != nil {
		utils.ErrorAesEncrypt.Message = err.Error()
		fail(c, utils.ErrorAesEncrypt)
		return
	}

	mm := model.ModelArts{
		Base:       model.Base{ID: 1},
		UserName:   dtoInput.UserName,
		AK:         dtoInput.AK,
		SK:         base64.StdEncoding.EncodeToString(encryptSK),
		ProjectId:  dtoInput.ProjectId,
		BucketName: dtoInput.BucketName,
		Site:       dtoInput.Site,
	}

	err = dao.ModelArtsCreate(&mm)
	if err != nil {
		utils.ErrModelArtsCreate.Message = err.Error()
		fail(c, utils.ErrModelArtsCreate)
		return
	}

	respWithData(c, map[string]interface{}{
		"id": mm.ID,
	})
}

func ModelArtsList(c *gin.Context) {
	dtoMa := dto.ModelArts{}

	ms, _, err := dao.ModelArtsList(&dtoMa)
	if err != nil {
		utils.ErrModelArtsList.Message = err.Error()
		fail(c, utils.ErrModelArtsList)
		return
	}

	if len(ms) == 0 {
		respWithData(c, model.ModelArts{})
		return
	}

	bytesPass, _ := base64.StdEncoding.DecodeString(ms[0].SK)
	decryptSK, err := utils.AesDecrypt(bytesPass, utils.DefaultKey)
	if err != nil {
		utils.ErrorAesDecrypt.Message = err.Error()
		fail(c, utils.ErrorAesDecrypt)
		return
	}

	ms[0].SK = string(decryptSK)
	respWithData(c, ms[0])

	return
}
