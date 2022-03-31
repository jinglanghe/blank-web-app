package model

import uuid "github.com/satori/go.uuid"

type AlertMetric struct {
	Base
	UUID  uuid.UUID `gorm:"type:uuid;unique" json:"uuid"`
	Type  string    `json:"type"`
	Count int64     `json:"count"`
}

type AlertMetricType struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}
