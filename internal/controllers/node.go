package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dto"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/service"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/utils"
	"strconv"
	"strings"
)

func registerNode(rg *gin.RouterGroup) {
	ctrl := &nodeController{}

	g := rg.Group("/nodes")
	g.GET("", ctrl.list)
	g.GET("/overview", ctrl.overview)
	g.GET("/pod-info", ctrl.podInfo)
	g.POST("/:id/label", ctrl.label)
	g.POST("/:id/taint", ctrl.taint)
}

var (
	nodeDao = &dao.Nodes{}
)

type nodeController struct {
	BaseController
}

func (n *nodeController) list(c *gin.Context) {
	var listDto dto.NodeList
	var resourceQuota *model.ResourceQuota
	var err error
	if !n.BindAndValidate(c, &listDto) {
		return
	}

	if listDto.ResourceQuotaID != 0 {
		listDto.PageSize = 0
		listDto.PageNum = 0
		resourceQuota, err = resourceQuotaDao.Get(listDto.ResourceQuotaID)
		if err != nil {
			failWithDBError(c, err)
			return
		}
		listDto.AvlCpuNum = resourceQuota.CpuNum
		listDto.AvlMem = resourceQuota.Mem
	}

	total, err := nodeDao.Count(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	nodes, err := nodeDao.List(&listDto)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	if listDto.ResourceQuotaID != 0 {
		nodes, total = n.filterQuota(nodes, resourceQuota)
	}

	respWithData(c, map[string]interface{}{"items": nodes, "total": total})
}

//return nodes that meet the resource quota
func (n *nodeController) filterQuota(nodes []model.Node, quota *model.ResourceQuota) ([]model.Node, int64) {
	var avlNum int64
	for i := 0; i < len(nodes); {
		avlNum = 0
		node := nodes[i]
		for _, dev := range node.Devs {
			if quota.EqualNodeDevice(dev) {
				avlNum = avlNum + dev.AvlNum
			}
		}
		if avlNum < quota.Num {
			nodes = append(nodes[:i], nodes[i+1:]...)
		} else {
			i++
		}
	}
	return nodes, int64(len(nodes))
}

func (n *nodeController) overview(c *gin.Context) {
	var nodeOverview model.NodeOverview
	nodes, err := nodeDao.ListWithDevs()
	if err != nil {
		failWithDBError(c, err)
		return
	}

	for _, node := range nodes {
		allocatedReses, err := service.NodeAllocatedReses(node.Name)
		if err != nil {
			utils.ErrNodeRequests.Message = err.Error()
			fail(c, utils.ErrNodeRequests)
			return
		}
		switch strings.ToLower(node.Type) {
		case model.CPU:
			if node.Status == model.NodeStatusReady {
				nodeOverview.ResStats.Cpu.Available++
				if n.Occupied(node, strings.ToUpper(model.CPU), allocatedReses) {
					nodeOverview.ResStats.Cpu.Occupied++
				}
			}
		case model.GPU:
			if node.Status == model.NodeStatusReady {
				nodeOverview.ResStats.Gpu.Available++
				if n.Occupied(node, strings.ToUpper(model.GPU), allocatedReses) {
					nodeOverview.ResStats.Gpu.Occupied++
				}
			}
		case model.NPU:
			if node.Status == model.NodeStatusReady {
				nodeOverview.ResStats.Npu.Available++
				if n.Occupied(node, strings.ToUpper(model.NPU), allocatedReses) {
					nodeOverview.ResStats.Npu.Occupied++
				}
			}
		}
	}

	respWithData(c, nodeOverview)
}

func (n *nodeController) Occupied(node model.Node, _type string, allocatedReses map[string]service.AllocatedRes) bool {
	for _, dev := range node.Devs {
		if dev.Type != _type {
			continue
		}
		for _, allocatedRes := range allocatedReses {
			if allocatedRes.EqualNodeDevice(&dev) && allocatedRes.Value() > 0 {
				return true
			}
		}
	}
	return false
}

func (n *nodeController) podInfo(c *gin.Context) {
	var listDto dto.NodePodList
	if !n.BindAndValidate(c, &listDto) {
		return
	}

	var nodesPodInfo []model.PodInfo

	ids, err := n.parseIds(listDto.Ids)
	if err != nil {
		utils.ErrorValidation.Message = err.Error()
		fail(c, utils.ErrorValidation)
		return
	}

	nodes, err := nodeDao.Gets(ids)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	for _, node := range nodes {
		podList, err := service.NodePodList(node.Name)
		if err != nil {
			utils.ErrNodePodsList.Message = err.Error()
			fail(c, utils.ErrNodePodsList)
			return
		}
		nodesPodInfo = append(nodesPodInfo, model.PodInfo{
			Id:     node.ID,
			PodNum: len(podList.Items),
		})
	}

	respWithData(c, nodesPodInfo)
}

func (n *nodeController) parseIds(idsStr string) ([]int64, error) {
	if len(idsStr) == 0 {
		return nil, nil
	}
	idStrArr := strings.Split(strings.TrimSpace(idsStr), "|")
	var ids []int64
	for _, idStr := range idStrArr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, int64(id))
	}
	return ids, nil
}

func (n *nodeController) label(c *gin.Context) {
	var createDto dto.NodeLabelCreate
	if !n.BindAndValidate(c, &createDto) {
		return
	}

	node, err := nodeDao.Get(createDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	err = service.NodeSetLabel(node, createDto.Key, createDto.Value)
	if err != nil {
		utils.ErrNodeLabelSet.Message = err.Error()
		fail(c, utils.ErrNodeLabelSet)
		return
	}

	ok(c)
}

func (n *nodeController) taint(c *gin.Context) {
	var createDto dto.NodeTaintCreate
	if !n.BindAndValidate(c, &createDto) {
		return
	}

	node, err := nodeDao.Get(createDto.ID)
	if err != nil {
		failWithDBError(c, err)
		return
	}

	err = service.NodeAddTaint(node, createDto.Key, createDto.Value, createDto.Effect)
	if err != nil {
		utils.ErrNodeTaintAdd.Message = err.Error()
		fail(c, utils.ErrNodeTaintAdd)
		return
	}

	ok(c)
}
