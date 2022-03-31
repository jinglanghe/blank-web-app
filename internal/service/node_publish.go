package service

import (
	"context"
	"encoding/json"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/broker"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const (
	Online       = "online"
	Offline      = "offline"
	UpdateStatus = "updateStatus"
)

var (
	apEdgeTopic = "apedgenode"
)

// NodeMsg action==updateStatus indicates the node status change:
// 	When the status is Ready, it means the node status is changed from NotReady to Ready,
// 	When the status is NotReady, the node status is changed from Ready to NotReady.
type NodeMsg struct {
	ID         int64              `json:"id"`
	Action     string             `json:"action"`
	Name       string             `json:"name"`
	Status     string             `json:"status"`
	Roles      string             `json:"roles"`
	Capacity   v1.ResourceList    `json:"capacity"`
	SystemInfo *v1.NodeSystemInfo `json:"systemInfo"`
}

func (n *NodeMsg) Marshal() []byte {
	bytes, _ := json.Marshal(n)
	return bytes
}

func NodePublishMsg(id int64, action string, name string) {
	nodeMsg := NodeMsg{
		ID:     id,
		Action: action,
		Name:   name,
	}
	if action == Online || action == UpdateStatus {
		n, err := clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			logging.Error(err).Msg("publish msg to mq failed: get node failed")
			return
		}
		nodeMsg.Status = getNodeStatus(n.Status.Conditions)
		nodeMsg.Roles = strings.Join(getNodeRoles(n), ",")
		nodeMsg.Capacity = n.Status.Capacity.DeepCopy()
		nodeMsg.SystemInfo = n.Status.NodeInfo.DeepCopy()
	}

	msg := broker.Message{
		Header: nil,
		Body:   nodeMsg.Marshal(),
	}
	err := mq.Publish(apEdgeTopic, &msg)
	if err != nil {
		logging.Error(err).Msg("publish msg to mq failed")
		return
	}
	logging.Debug().Msgf("publish msg success: node.name:%v, action:%v", name, action)
}
