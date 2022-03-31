package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dto"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/service"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/service/aaa"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/utils"
)

func registerOrgResource(rg *gin.RouterGroup) {
	ctrl := &orgResourceController{}

	g := rg.Group("/org-resources")
	g.PATCH("/:orgId", ctrl.create)
	g.GET("", ctrl.list)
	g.GET("/:id", ctrl.get)
	g.DELETE("/:id", ctrl.delete)
}

var (
	orgResourceDao = &dao.OrgResource{}
)

type orgResourceController struct {
	BaseController
}

func (o *orgResourceController) create(c *gin.Context) {
	var createDto dto.OrgResourceCreate

	if !o.BindAndValidate(c, &createDto) {
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

	org, err := aaa.OrgDetail(token, createDto.OrgId)
	if err != nil {
		utils.ErrOrgDetail.Message = err.Error()
		fail(c, utils.ErrOrgDetail)
		return
	}
	if org.ID == 0 {
		fail(c, utils.ErrOrgDetail)
		return
	}

	createDto.Duration = 0
	createDto.Source = 0
	createDto.CreatorID = u.ID
	createDto.CreatorName = u.Username
	createDto.OrgName = org.Account

	newNodes, err := nodeDao.Gets(createDto.NodeIds)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	orgResource, err := orgResourceDao.GetByOrgId(createDto.OrgId)
	if err != nil {
		failWithDBError(c, err)
		return
	}
	if orgResource.ID == 0 {
		newOrgResource := model.OrgResource{
			OrgID:       createDto.OrgId,
			OrgName:     createDto.OrgName,
			Duration:    createDto.Duration,
			Source:      createDto.Source,
			CreatorID:   createDto.CreatorID,
			CreatorName: createDto.CreatorName,
			Nodes:       newNodes,
		}
		for _, newNode := range newNodes {
			err := service.NodeCordonOrUncordon(newNode.Name, false)
			if err != nil {
				fail(c, utils.ErrNodeUnCordon)
				return
			}
		}
		err := orgResourceDao.Create(&newOrgResource)
		if err != nil {
			failWithDBError(c, err)
			return
		}
	} else {
		for _, oldNode := range orgResource.Nodes {
			if !orgResourceDao.Exist(oldNode, createDto.NodeIds) {
				err := service.NodeCordonOrUncordon(oldNode.Name, true)
				if err != nil {
					fail(c, utils.ErrNodeCordon)
					return
				}
				if err := nodeDao.FreeNode(&oldNode); err != nil {
					failWithDBError(c, err)
					return
				}
			}
		}
		for _, newNode := range newNodes {
			isNewNode := true
			for _, oldNode := range orgResource.Nodes {
				if oldNode.ID == newNode.ID {
					isNewNode = false
					continue
				}
			}

			if isNewNode {
				err := service.NodeCordonOrUncordon(newNode.Name, false)
				if err != nil {
					fail(c, utils.ErrNodeUnCordon)
					return
				}
			}
		}
		orgResource.Nodes = newNodes
		if err := orgResourceDao.Save(orgResource); err != nil {
			failWithDBError(c, err)
			return
		}
	}

	service.RefreshNodeLabel()

	ok(c)
}

func (o *orgResourceController) list(c *gin.Context) {
	var listDto dto.OrgResourceList
	if !o.BindAndValidate(c, &listDto) {
		return
	}

	total, err := orgResourceDao.Count(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	orgResources, err := orgResourceDao.List(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	respWithData(c, map[string]interface{}{"items": orgResources, "total": total})
}

func (o *orgResourceController) get(c *gin.Context) {
	var getDto dto.OrgResourceGet
	if !o.BindAndValidate(c, &getDto) {
		return
	}

	orgResource, err := orgResourceDao.Get(getDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}
	respWithData(c, orgResource)
}

func (o *orgResourceController) delete(c *gin.Context) {
	var deleteDto dto.OrgResourceDelete
	if !o.BindAndValidate(c, &deleteDto) {
		return
	}

	orgResource, err := orgResourceDao.Get(deleteDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	if err := orgResourceDao.Delete(orgResource); err != nil {
		failWithDBError(c, err)
		return
	}

	ok(c)
}
