package service

import (
	"context"
	"fmt"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/utils"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
	alertconfig "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	KubeSystemNS   = "kube-system"
	PromConfigmap  = "prometheus-server"
	AlertConfigmap = "prometheus-alertmanager"
	AlertConfigKey = "alertmanager.yml"

	fnFmt = map[model.RuleType]string{
		model.ServiceRuleType:  "%d_service.rules",
		model.ResourceRuleType: "%d_resource.rules",
	}
)

func GeneratePromRule(rule model.BaseRule) error {
	fn := fmt.Sprintf(fnFmt[rule.GetRuleType()], rule.GetOrgId())

	var groups *rulefmt.RuleGroups
	ruleCm, err := loadConfigmap(PromConfigmap)
	if err != nil {
		return err
	}
	if _, exist := ruleCm.Data[fn]; !exist {
		groups = &rulefmt.RuleGroups{Groups: []rulefmt.RuleGroup{}}
	} else {
		_groups, errs := rulefmt.Parse([]byte(ruleCm.Data[fn]))
		if len(errs) > 0 {
			for _, err := range errs {
				logging.Error(err).Str("file", fn).Msg("cannot parse rule file")
			}
			return fmt.Errorf("cannot parse rule file: %s", fn)
		}
		groups = _groups
	}

	newGroup, err := rule.GeneratePromRuleGroup()
	if err != nil {
		logging.Error(err).Msg("create prometheus group rule error")
		return err
	}

	found := false
	for idx, group := range groups.Groups {
		if group.Name == rule.GetUUID() {
			groups.Groups[idx] = *newGroup
			found = true
		}
	}
	if found == false {
		groups.Groups = append(groups.Groups, *newGroup)
	}

	if err := saveConfigmap(ruleCm, fn, groups); err != nil {
		logging.Error(err).Str("file", fn).Msg("write rule data to file error")
		return err
	}

	if err := newAlertRouteEntry(rule); err != nil {
		logging.Error(err).Msg("generate alert manager error")
		return err
	}

	return nil
}

func UpdatePromRule(rule model.BaseRule) error {
	fn := fmt.Sprintf(fnFmt[rule.GetRuleType()], rule.GetOrgId())

	ruleCm, err := loadConfigmap(PromConfigmap)
	if err != nil {
		return err
	}

	if _, exist := ruleCm.Data[fn]; !exist {
		logging.Error(err).Str("file", fn).Msg("file not exist")
		return err
	}

	groups, errs := rulefmt.Parse([]byte(ruleCm.Data[fn]))
	if len(errs) > 0 {
		for _, err := range errs {
			logging.Error(err).Str("file", fn).Msg("cannot parse rule file")
		}
		return fmt.Errorf("cannot parse rule file: %s", fn)
	}

	for idx, group := range groups.Groups {
		if group.Name != rule.GetUUID() {
			continue
		}

		newGroup, err := rule.GeneratePromRuleGroup()
		if err != nil {
			logging.Error(err).Msg("create prometheus group rule error")
			return err
		}
		groups.Groups[idx] = *newGroup

		if err := saveConfigmap(ruleCm, fn, groups); err != nil {
			logging.Error(err).Str("file", fn).Msg("write rule data to file error")
			return err
		}

		if err := updateAlertConfig(rule); err != nil {
			logging.Error(err).Msg("update alert manager error")
			return err
		}

		return nil
	}

	return fmt.Errorf("not found group: %s", rule.GetUUID())
}

func DeletePromRule(rule model.BaseRule) *utils.CodeMessage {
	fn := fmt.Sprintf(fnFmt[rule.GetRuleType()], rule.GetOrgId())

	ruleCm, err := loadConfigmap(PromConfigmap)
	if err != nil {
		utils.ErrorConfigmapOp.Message = err.Error()
		return utils.ErrorConfigmapOp
	}

	if _, exist := ruleCm.Data[fn]; !exist {
		logging.Error(err).Str("file", fn).Msg("file not exist")
		utils.ErrorConfigmap404.Message = fmt.Sprintf("%s %s", fn, "file not exist")
		return utils.ErrorConfigmap404
	}

	groups, errs := rulefmt.Parse([]byte(ruleCm.Data[fn]))
	if len(errs) > 0 {
		for _, err := range errs {
			logging.Error(err).Str("file", fn).Msg("cannot parse rule file")
		}
		utils.ErrorRuleValidate.Message = fmt.Sprintf("cannot parse rule file: %s", fn)
		return utils.ErrorRuleValidate
	}

	for i := 0; i < len(groups.Groups); {
		if groups.Groups[i].Name == rule.GetUUID() {
			groups.Groups = append(groups.Groups[:i], groups.Groups[i+1:]...)

			if err := deleteAlterConfig(rule.GetUUID()); err != nil {
				return err
			}

			if err := saveConfigmap(ruleCm, fn, groups); err != nil {
				logging.Error(err).Str("file", fn).Msg("write rule data to file error")
				utils.ErrorConfigmapOp.Message = err.Error()
				return utils.ErrorConfigmapOp
			}

			logging.Debug().Str("result", "success").
				Str("uuid", rule.GetUUID()).Msg("delete prometheus rule")
		} else {
			i++
		}
	}

	logging.Debug().Str("action", "delete prometheus rule").
		Str("uuid", rule.GetUUID()).
		Msg("not found group")
	return nil
}

func newAlertRouteEntry(rule model.BaseRule) error {
	cm, err := loadConfigmap(AlertConfigmap)
	if err != nil {
		return err
	}
	amConfig, err := alertconfig.Load(cm.Data[AlertConfigKey])
	if err != nil {
		return err
	}

	newRoute := rule.GenerateAlertRoute()
	amConfig.Route.Routes = append(amConfig.Route.Routes, newRoute)

	if err := saveConfigmap(cm, AlertConfigKey, amConfig); err != nil {
		return err
	}

	return nil
}

func updateAlertConfig(rule model.BaseRule) error {
	cm, err := loadConfigmap(AlertConfigmap)
	if err != nil {
		return err
	}
	amConfig, err := alertconfig.Load(cm.Data[AlertConfigKey])
	if err != nil {
		return err
	}

	for idx, route := range amConfig.Route.Routes {
		if route.Match["uuid"] != rule.GetUUID() {
			continue
		}

		amConfig.Route.Routes[idx] = rule.GenerateAlertRoute()
		if err := saveConfigmap(cm, AlertConfigKey, amConfig); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("not found route: %s", rule.GetUUID())
}

func deleteAlterConfig(uuid string) *utils.CodeMessage {
	cm, err := loadConfigmap(AlertConfigmap)
	if err != nil {
		utils.ErrorConfigmapOp.Message = err.Error()
		return utils.ErrorConfigmapOp
	}
	amConfig, err := alertconfig.Load(cm.Data[AlertConfigKey])
	if err != nil {
		utils.ErrorConfigmapOp.Message = err.Error()
		return utils.ErrorConfigmapOp
	}

	for idx, route := range amConfig.Route.Routes {
		if uuid != route.Match["uuid"] {
			continue
		}

		amConfig.Route.Routes = append(amConfig.Route.Routes[:idx], amConfig.Route.Routes[idx+1:]...)
		if err := saveConfigmap(cm, AlertConfigKey, amConfig); err != nil {
			utils.ErrorConfigmapOp.Message = err.Error()
			return utils.ErrorConfigmapOp
		}

		return nil
	}

	utils.ErrorConfigmap404.Message = fmt.Sprintf("not found route: %s", uuid)
	return utils.ErrorConfigmap404
}

func saveConfigmap(cm *v1.ConfigMap, fn string, data interface{}) error {
	d, err := yaml.Marshal(data)
	if err != nil {
		logging.Error(err).Msg("marshal rule data to byte")
		return err
	}

	cm.Data[fn] = string(d)
	if _, err := clientset.CoreV1().ConfigMaps(KubeSystemNS).Update(context.Background(), cm, metav1.UpdateOptions{}); err != nil {
		logging.Error(err).Str("file", fn).Msg("write data to configmap error")
		return err
	}

	return nil
}

func loadConfigmap(name string) (*v1.ConfigMap, error) {
	cm, err := clientset.CoreV1().ConfigMaps(KubeSystemNS).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return cm, nil
}
