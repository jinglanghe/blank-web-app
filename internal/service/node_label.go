package service

import (
	"context"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const NodeOrganizationLabelKey = "apulis.com/organization.name"

func RefreshNodeLabel() {
	nodes, err := nodeDao.ListAll()
	if err != nil {
		logging.Error(err).Msg("get nodes error")
		return
	}
	if err = nodesSetLabel(nodes); err != nil {
		logging.Error(err).Msg("refresh nodes label error")
	}
}

func nodesSetLabel(nodes []model.Node) error {
	for _, node := range nodes {
		err := NodeSetLabel(&node, NodeOrganizationLabelKey, node.OrgResource.OrgName)
		if err != nil {
			return err
		}
	}
	return nil
}

func NodeSetLabel(node *model.Node, key string, value string) error {
	n, err := clientset.CoreV1().Nodes().Get(context.TODO(), node.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	labels := n.GetLabels()
	labels[key] = value
	n.SetLabels(labels)

	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), n, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
