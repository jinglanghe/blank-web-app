package service

import (
	"context"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"github.com/apulis/sdk/go-utils/logging"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
)

func RefreshQuota() {
	resources, err := userGroupResourceDao.ListALl()
	if err != nil {
		logging.Error(err).Msg("get user group resources error")
		return
	}

	nodeDeviceTypeOrgs, err := nodeDeviceDao.OrgDevises()
	if err != nil {
		logging.Error(err).Msg("get node devises error")
		return
	}

	orgAvlMemes, err := nodeDao.OrgAvlMemes()
	if err != nil {
		logging.Error(err).Msg("get org memes error")
		return
	}

	listReses := userGroupResourceDao.MatchResources(resources, nodeDeviceTypeOrgs, orgAvlMemes)
	for _, listRes := range listReses {
		_ = setQuota(clientset, &listRes)
	}
}

func setQuota(c *kubernetes.Clientset, resourceList *model.UserGroupResourceList) error {
	resourceQuota := &v1.ResourceQuota{ObjectMeta: metav1.ObjectMeta{Name: resourceList.QuotaName()}}
	resourceQuota.Spec.Hard = make(map[v1.ResourceName]resource.Quantity, 0)

	if resourceList.MemMin > 0 {
		memQuantity, err := resource.ParseQuantity(strconv.Itoa(int(resourceList.MemMin)))
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}
		resourceQuota.Spec.Hard["requests.memory"] = memQuantity
	}

	if resourceList.MemMax > 0 {
		memQuantity, err := resource.ParseQuantity(strconv.Itoa(int(resourceList.MemMax)))
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}
		resourceQuota.Spec.Hard["limits.memory"] = memQuantity
	}

	for _, quota := range resourceList.Quotas {
		if quota.Type == strings.ToUpper(model.CPU) {
			cpuQuantity, err := resource.ParseQuantity(strconv.Itoa(int(quota.Min)))
			if err != nil {
				logging.Error(err).Msg("")
				return err
			}
			resourceQuota.Spec.Hard["requests.cpu"] = cpuQuantity

			cpuQuantity, err = resource.ParseQuantity(strconv.Itoa(int(quota.Max)))
			if err != nil {
				logging.Error(err).Msg("")
				return err
			}
			resourceQuota.Spec.Hard["limits.cpu"] = cpuQuantity
			continue
		}

		quantity, err := resource.ParseQuantity(strconv.Itoa(int(quota.Min)))
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}
		resourceQuota.Spec.Hard[v1.ResourceName("requests."+quota.Model)] = quantity

		quantity, err = resource.ParseQuantity(strconv.Itoa(int(quota.Max)))
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}
		resourceQuota.Spec.Hard[v1.ResourceName("limits."+quota.Model)] = quantity
	}

	_, err := c.CoreV1().ResourceQuotas(resourceList.Namespace()).
		Get(context.TODO(), resourceQuota.Name, metav1.GetOptions{})
	if err != nil {
		_, err = c.CoreV1().ResourceQuotas(resourceList.Namespace()).
			Create(context.TODO(), resourceQuota, metav1.CreateOptions{})
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}
	} else {
		_, err = c.CoreV1().ResourceQuotas(resourceList.Namespace()).
			Update(context.TODO(), resourceQuota, metav1.UpdateOptions{})
		if err != nil {
			logging.Error(err).Msg("")
			return err
		}
	}
	logging.Info().Msgf("refresh resource quota: %v", resourceList.OrgName)
	return nil
}
