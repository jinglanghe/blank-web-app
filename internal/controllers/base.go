package controllers

import (
	"errors"
	"fmt"
	//"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/logging"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/utils"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseController struct{}

func (b *BaseController) BindAndValidate(c *gin.Context, params interface{}) bool {
	if err := bind(c, params); err != nil {
		failValidate(c, err.Error())
		return false
	}

	return true
}

func bindAndValidate(c *gin.Context, params interface{}) bool {
	if err := bind(c, params); err != nil {
		failValidate(c, err.Error())
		return false
	}

	return true
}

func bind(c *gin.Context, params interface{}) error {
	_ = c.ShouldBindUri(params)
	if err := c.ShouldBind(params); err != nil {
		if fieldErr, ok := err.(validator.ValidationErrors); ok {
			var tagErrorMsg []string
			for _, v := range fieldErr {
				if _, has := utils.ValidateErrorMessage[v.Tag()]; has {
					tagErrorMsg = append(tagErrorMsg, fmt.Sprintf(utils.ValidateErrorMessage[v.Tag()], v.Field(), v.Value()))
				} else {
					tagErrorMsg = append(tagErrorMsg, err.Error())
				}
			}

			logging.Error(err).Msg("")
			return errors.New(strings.Join(tagErrorMsg, ","))
		}
	}

	return nil
}

func failValidate(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"code": utils.ErrorValidation.Code,
		"msg":  msg,
		"data": struct{}{},
	})
}

func ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    utils.ActionSuccess.Code,
		"message": utils.ActionSuccess.Message,
		"data":    struct{}{},
	})
}

func respWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": utils.ActionSuccess.Code,
		"msg":  utils.ActionSuccess.Message,
		"data": data,
	})
}

func respResourceNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": utils.ErrorNotFound.Code,
		"msg":  utils.ErrorNotFound.Message,
		"data": struct{}{},
	})
}

func fail(c *gin.Context, err *utils.CodeMessage) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": err.Code,
		"msg":  err.Message,
		"data": struct{}{},
	})
}

func failWithStatus(c *gin.Context, status int, err *utils.CodeMessage) {
	c.JSON(status, gin.H{
		"code": err.Code,
		"msg":  err.Message,
		"data": struct{}{},
	})
}

func failWithData(c *gin.Context, err *utils.CodeMessage, data interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": err.Code,
		"msg":  err.Message,
		"data": data,
	})
}

func failWithHttpCode(c *gin.Context, httpCode int, err *utils.CodeMessage) {
	c.JSON(httpCode, gin.H{
		"code": err.Code,
		"msg":  err.Message,
	})
}

func failWithDBError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ErrorRecordNotExist.Message = err.Error()
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrorRecordNotExist.Code,
			"msg":  utils.ErrorRecordNotExist.Message,
		})
	} else {
		utils.ErrorDatabaseOp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ErrorDatabaseOp.Code,
			"msg":  utils.ErrorDatabaseOp.Message,
		})
	}
}
