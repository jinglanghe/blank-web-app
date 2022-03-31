package service

import (
	"context"
	"github.com/apulis/sdk/go-utils/logging"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNamespace(name string) error {
	err := initClientSet()
	if err != nil {
		return err
	}

	ns := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		logging.Error(err).Msg("")
		return err
	}

	logging.Info().Msgf("create namespace: %v", name)

	return nil
}

func NamespacePosList(name string) (*v1.PodList, error) {
	return clientset.CoreV1().Pods(name).List(context.TODO(), metav1.ListOptions{})
}

func NamespaceAllocatedReses(name string) (allocatedReses map[string]AllocatedRes, err error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podList, err := NamespacePosList(name)
	if err != nil {
		return
	}

	_, allocatedReses = getPodsTotalRequestsAndLimits(podList, nodes.Items)

	return
}
