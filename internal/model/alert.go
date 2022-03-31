package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

const (
	StatusToBeResolved = 0
	StatusResolved     = 1
	StatusUnResolved   = 2
)

type ServiceAlert struct {
	Base
	UUID       uuid.UUID      `gorm:"type:uuid;unique" json:"uuid"`
	OrgId      int64          `gorm:"column:org_id" json:"orgId"`
	Name       string         `gorm:"name" json:"name"`
	DataFilter datatypes.JSON `gorm:"data_filter" json:"dataFilter"`
	Condition  datatypes.JSON `json:"condition"`
	Rule       string         `json:"rule"`
	Receivers  datatypes.JSON `json:"receivers"`
	Labels     datatypes.JSON `json:"labels"`
	Type       string         `json:"type"`
	Source     string         `json:"source"`
	Status     int            `json:"status"`
}

func (ServiceAlert) TableName() string {
	return "service_alerts"
}

type ResourceAlert struct {
	Base
	UUID       uuid.UUID      `gorm:"type:uuid;unique" json:"uuid"`
	OrgId      int64          `gorm:"column:org_id" json:"orgId"`
	Name       string         `gorm:"name" json:"name"`
	Level      string         `json:"level"`
	DataFilter datatypes.JSON `gorm:"data_filter" json:"dataFilter"`
	Condition  datatypes.JSON `json:"condition"`
	Rule       string         `json:"rule"`
	Receivers  datatypes.JSON `json:"receivers"`
	Labels     datatypes.JSON `json:"labels"`
	Type       string         `json:"type"`
	Source     string         `json:"source"`
	Status     int            `json:"status"`
}

func (ResourceAlert) TableName() string {
	return "resource_alerts"
}
