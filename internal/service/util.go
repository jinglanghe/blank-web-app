package service

import (
	"github.com/apulis/bmod/aistudio-aom/internal/model"
	"github.com/apulis/bmod/aistudio-aom/internal/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
)

type AllocatedRes struct {
	Arch        string            `json:"arch"`
	Model       string            `json:"model"`
	ComputeType string            `json:"computeType"`
	Series      string            `json:"series"`
	Quantity    resource.Quantity `json:"_"`
}

func (a *AllocatedRes) Key() string {
	return a.Arch + a.Model + a.ComputeType + a.Series
}

func (a *AllocatedRes) Value() int64 {
	return utils.QuantityRoundUp(&a.Quantity)
}

func (a *AllocatedRes) EqualNodeDevice(d *model.NodeDevice) bool {
	if d.IsCPU() {
		if a.Model == string(corev1.ResourceCPU) && a.Arch == d.Arch {
			return true
		}
	} else {
		if a.Model == d.Model && a.ComputeType == d.ComputeType && a.Series == d.Series {
			return true
		}
	}
	return false
}

func (a *AllocatedRes) EqualQuota(q *model.Quota) bool {
	if q.IsCPU() {
		if a.Model == string(corev1.ResourceCPU) && a.Arch == q.Arch {
			return true
		}
	} else {
		if a.Model == q.Model && a.ComputeType == q.ComputeType && a.Series == q.Series {
			return true
		}
	}
	return false
}

func (a *AllocatedRes) IsMem() bool {
	if a.Model == string(corev1.ResourceMemory) {
		return true
	}
	return false
}

func getPodsTotalRequestsAndLimits(podList *corev1.PodList, nodes []corev1.Node) (reqs map[string]AllocatedRes, limits map[string]AllocatedRes) {
	reqs, limits = map[string]AllocatedRes{}, map[string]AllocatedRes{}
	for _, pod := range podList.Items {
		node := getNode(nodes, pod.Spec.NodeName)
		podReqs, podLimits := resourcehelper.PodRequestsAndLimits(&pod)
		for podReqName, podReqValue := range podReqs {
			allocatedRes := AllocatedRes{
				Arch:        node.Status.NodeInfo.Architecture,
				Model:       string(podReqName),
				ComputeType: getComputeType(node.ObjectMeta.Labels),
				Series:      getSeries(node.ObjectMeta.Labels),
				Quantity:    podReqValue.DeepCopy(),
			}
			if value, ok := reqs[allocatedRes.Key()]; !ok {
				reqs[allocatedRes.Key()] = allocatedRes
			} else {
				value.Quantity.Add(podReqValue)
				reqs[allocatedRes.Key()] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			allocatedRes := AllocatedRes{
				Arch:        node.Status.NodeInfo.Architecture,
				Model:       string(podLimitName),
				ComputeType: getComputeType(node.ObjectMeta.Labels),
				Series:      getSeries(node.ObjectMeta.Labels),
				Quantity:    podLimitValue.DeepCopy(),
			}
			if value, ok := limits[allocatedRes.Key()]; !ok {
				limits[allocatedRes.Key()] = allocatedRes
			} else {
				value.Quantity.Add(podLimitValue)
				limits[allocatedRes.Key()] = value
			}
		}
	}
	return
}
