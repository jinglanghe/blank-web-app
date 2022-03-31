package model

import (
	alertconfig "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	//"github.com/prometheus/prometheus/pkg/rulefmt"
)

type RuleType int

const (
	_                RuleType = iota // 0
	ServiceRuleType                  // 1
	ResourceRuleType                 // 2
	StatusOff        = 0
	StatusOn         = 1
)

var (
	defaultReceiver = "default-receiver"
)

type BaseRule interface {
	GetUUID() string
	GetRuleType() RuleType
	GetOrgId() int64
	GeneratePromRuleGroup() (*rulefmt.RuleGroup, error)
	GenerateAlertRoute() *alertconfig.Route
}
