package dto

import (
	"gorm.io/datatypes"
)

type ServiceAlertListDto struct {
	BaseListDto
	StartTime int64          `form:"startTime" json:"startTime" time_format:"unix" binding:"required"`
	EndTime   int64          `form:"endTime" json:"endTime" time_format:"unix" binding:"required"`
	Name      string         `form:"name" json:"name"`
	Labels    datatypes.JSON `json:"labels"`
	Type      string         `json:"type"`
	Source    string         `json:"source"`
}

type ServiceAlertsDeleteDto struct {
	UUIDs []string `uri:"uuids" json:"uuids" binding:"required"`
}

type ServiceAlertStatusDto struct {
	UUID   string `uri:"uuid" json:"uuid" binding:"required"`
	Status *int   `uri:"status" json:"status" binding:"required"`
}

type ResourceAlertListDto struct {
	BaseListDto
	Level     int            `form:"level" json:"level"`
	StartTime int64          `form:"startTime" json:"startTime" time_format:"unix" binding:"required"`
	EndTime   int64          `form:"endTime" json:"endTime" time_format:"unix" binding:"required"`
	Name      string         `form:"name" json:"name"`
	Labels    datatypes.JSON `json:"labels"`
	Type      string         `json:"type"`
	Source    string         `json:"source"`
}

type ResourceAlertsDeleteDto struct {
	UUIDs []string `uri:"uuids" json:"uuids" binding:"required"`
}

type ResourceAlertStatusDto struct {
	UUID   string `uri:"uuid" json:"uuid" binding:"required"`
	Status *int   `uri:"status" json:"status" binding:"required"`
}
