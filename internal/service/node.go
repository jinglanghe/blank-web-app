package service

import (
	"context"
	"fmt"
	"github.com/apulis/bmod/aistudio-aom/internal/dao"
	"github.com/apulis/bmod/aistudio-aom/internal/model"
	"github.com/apulis/sdk/go-utils/logging"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/drain"
	"strings"
)

func RefreshModelNode() {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logging.Error(err)
		return
	}

	var nodeItems []model.Node
	for _, node := range nodes.Items {
		nodeItem := model.Node{
			MachineID: getNodeUniqueId(&node),
			Name:      node.ObjectMeta.Name,
			Type:      strings.ToUpper(model.CPU),
			CpuArch:   node.Status.NodeInfo.Architecture,
			Role:      getNodeRole(node.ObjectMeta.Labels),
			Status:    getNodeStatus(node.Status.Conditions),
		}

		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				nodeItem.InternalIP = addr.Address
			}
		}
		if len(nodeItem.InternalIP) == 0 {
			nodeItem.InternalIP = nodeItem.Name
		}

		cpuNum, _ := node.Status.Capacity.Cpu().AsInt64()
		nodeItem.CpuNum = float32(cpuNum)
		avlCpuNum, _ := node.Status.Allocatable.Cpu().AsInt64()
		nodeItem.AvlCpuNum = float32(avlCpuNum)
		nodeItem.Mem, _ = node.Status.Capacity.Memory().AsInt64()
		nodeItem.AvlMem, _ = node.Status.Allocatable.Memory().AsInt64()

		series := getSeries(node.ObjectMeta.Labels)
		computeType := getComputeType(node.ObjectMeta.Labels)

		existNpu := false
		for k, v := range node.Status.Capacity {
			avlNum := getAvlNum(node.Status.Allocatable, k)
			if strings.Contains(string(k), model.GPU) {
				num, _ := v.AsInt64()
				nodeDevice := model.NodeDevice{
					Type:        strings.ToUpper(model.GPU),
					Model:       string(k),
					ComputeType: computeType,
					Series:      series,
					Num:         num,
					AvlNum:      avlNum,
				}
				if existNpu == false {
					nodeItem.Type = strings.ToUpper(model.GPU)
				}
				nodeItem.Devs = append(nodeItem.Devs, nodeDevice)
				continue
			}

			if strings.Contains(string(k), model.NPU) {
				num, _ := v.AsInt64()
				nodeDevice := model.NodeDevice{
					Type:        strings.ToUpper(model.NPU),
					Model:       string(k),
					ComputeType: computeType,
					Series:      series,
					Num:         num,
					AvlNum:      avlNum,
				}

				existNpu = true
				nodeItem.Type = strings.ToUpper(model.NPU)
				nodeItem.Devs = append(nodeItem.Devs, nodeDevice)
				continue
			}

			if strings.Contains(string(k), model.CPU) {
				num, _ := v.AsInt64()
				nodeDevice := model.NodeDevice{
					Type:   strings.ToUpper(model.CPU),
					Arch:   node.Status.NodeInfo.Architecture,
					Num:    num,
					AvlNum: avlNum,
				}
				nodeItem.Devs = append(nodeItem.Devs, nodeDevice)
			}
		}

		nodeItems = append(nodeItems, nodeItem)
	}

	if len(nodeItems) == 0 {
		return
	}

	err = NodeUpsert(nodeItems)
	if err != nil {
		logging.Error(err).Msg("NodeUpsert failed")
		return
	}
	return
}

func NodeUpsert(newNodes []model.Node) error {
	defaultOrgResource := dao.DefaultOrgResource()

	var newAddNodes []model.Node
	var deleteNodes []model.Node
	var deleteNodeDevIds []int64

	oldNodes, err := nodeDao.ListWithDevs()
	if err != nil {
		logging.Error(err).Msg("node upsert failed at get old nodes")
		return err
	}

	for _, oldNode := range oldNodes {
		newNodeIndex := nodeDao.FindNode(&oldNode, newNodes)
		if newNodeIndex > len(newNodes) {
			deleteNodes = append(deleteNodes, oldNode)
			continue
		}
		newNodes[newNodeIndex].ID = oldNode.ID
		newNodes[newNodeIndex].CreatedAt = oldNode.CreatedAt
		newNodes[newNodeIndex].OrgResourceId = oldNode.OrgResourceId
		for _, oldNodeDev := range oldNode.Devs {
			newNodeDevIndex := nodeDao.FindNodeDevice(&oldNodeDev, newNodes[newNodeIndex].Devs)
			if newNodeDevIndex > len(newNodes[newNodeIndex].Devs) {
				deleteNodeDevIds = append(deleteNodeDevIds, oldNodeDev.ID)
				continue
			}
			newNodes[newNodeIndex].Devs[newNodeDevIndex].CreatedAt = oldNodeDev.CreatedAt
			newNodes[newNodeIndex].Devs[newNodeDevIndex].ID = oldNodeDev.ID
		}

		if newNodes[newNodeIndex].Status != oldNode.Status {
			NodePublishMsg(newNodes[newNodeIndex].ID, UpdateStatus, newNodes[newNodeIndex].Name)
		}
	}

	//assign new node to the default organization
	for newNodeIndex, newNode := range newNodes {
		if newNode.ID == 0 {
			newNodes[newNodeIndex].OrgResourceId = defaultOrgResource.ID
			newAddNodes = append(newAddNodes, newNode)
		}
	}

	for _, deleteNode := range deleteNodes {
		err := nodeDao.Delete(&deleteNode)
		if err != nil {
			logging.Error(err).Msg(fmt.Sprintf("node upsert failed at delete node: %v", deleteNode.ID))
			return err
		}
		NodePublishMsg(deleteNode.ID, Offline, deleteNode.Name)
	}

	if len(deleteNodeDevIds) > 0 {
		err := nodeDeviceDao.Deletes(deleteNodeDevIds)
		if err != nil {
			logging.Error(err).Msg("node upsert failed at delete node devs")
			return err
		}
	}

	err = nodeDao.Save(newNodes)
	if err != nil {
		logging.Error(err).Msg("node upsert failed")
		return err
	}

	for _, newAddNode := range newAddNodes {
		newNodeIndex := nodeDao.FindNode(&newAddNode, newNodes)
		if newNodeIndex > len(newNodes) {
			continue
		}
		NodePublishMsg(newNodes[newNodeIndex].ID, Online, newNodes[newNodeIndex].Name)
	}

	logging.Info().Msgf("upsert %v nodes", len(newNodes))
	return nil
}

//In order to be compatible with the problem that the edge node cannot obtain the machine ID
func getNodeUniqueId(node *v1.Node) string {
	if node.Status.NodeInfo.MachineID != "" {
		return node.Status.NodeInfo.MachineID
	} else {
		return node.Name
	}
}

func getNodeRole(labels map[string]string) string {
	if role, exist := labels["kubernetes.io/role"]; exist {
		return role
	}
	if _, exist := labels["node-role.kubernetes.io/master"]; exist {
		return "master"
	}
	return ""
}

func getNodeRoles(node *v1.Node) []string {
	var roles []string
	labels := node.GetLabels()

	for labelKey, labelValue := range labels {
		if labelKey == "kubernetes.io/role" {
			roles = append(roles, labelValue)
		} else if strings.HasPrefix(labelKey, "node-role.kubernetes.io/") {
			role := strings.TrimPrefix(labelKey, "node-role.kubernetes.io/")
			roles = append(roles, role)
		}
	}
	return roles
}

func getNodeStatus(conditions []v1.NodeCondition) string {
	for _, c := range conditions {
		if c.Type == v1.NodeReady && c.Status == v1.ConditionTrue {
			return model.NodeStatusReady
		}
	}
	return model.NodeStatusNotReady
}

func getSeries(m map[string]string) string {
	for k, v := range m {
		if k == model.SERIES {
			return v
		}
	}

	return ""
}

func getComputeType(m map[string]string) string {
	for k, v := range m {
		if k == model.ComputeType {
			return v
		}
	}

	return ""
}

func getAvlNum(list v1.ResourceList, k v1.ResourceName) int64 {
	if v, exist := list[k]; exist {
		avlNum, _ := v.AsInt64()
		return avlNum
	}
	return 0
}

func NodeCordonOrUncordon(name string, desired bool) error {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logging.Error(err)
		return err
	}

	helper := &drain.Helper{
		Ctx:    context.TODO(),
		Client: clientset,
	}

	err = drain.RunCordonOrUncordon(helper, node, desired)
	return err
}

func NodePodList(name string) (*v1.PodList, error) {
	return clientset.CoreV1().Pods("").List(context.TODO(),
		metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + name,
		})
}

func NodeAllocatedReses(name string) (allocatedReses map[string]AllocatedRes, err error) {
	podList, err := NodePodList(name)
	if err != nil {
		return
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}

	_, allocatedReses = getPodsTotalRequestsAndLimits(podList, nodes.Items)
	return
}

func getNode(nodes []v1.Node, name string) v1.Node {
	for _, node := range nodes {
		if node.Name == name {
			return node
		}
	}
	return v1.Node{}
}

func NodeAddTaint(node *model.Node, key string, value string, effect string) error {
	n, err := clientset.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	n.Spec.Taints = append(n.Spec.Taints, v1.Taint{
		Key:    key,
		Value:  value,
		Effect: v1.TaintEffect(effect),
	})

	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), n, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
