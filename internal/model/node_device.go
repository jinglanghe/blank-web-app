package model

import (
	"gorm.io/gorm"
	"strings"
)

type NodeDevice struct {
	Base
	NodeID      int64  `json:"nodeId"`
	Key         string `json:"key"`
	Type        string `json:"type"`  //CPU|GPU|NPU
	Arch        string `json:"arch"`  //x86|arm64
	Model       string `json:"model"` //npu.huawei.com/NPU
	ComputeType string `json:"computeType"`
	Series      string `json:"series"` //a310
	Num         int64  `json:"num"`
	AvlNum      int64  `json:"avlNum"`
}

type NodeDeviceType struct {
	Key         string `json:"key"`
	Type        string `json:"type"`  //CPU|GPU|NPU
	Arch        string `json:"arch"`  //x86|arm64
	Model       string `json:"model"` //npu.huawei.com/NPU
	ComputeType string `json:"computeType"`
	Series      string `json:"series"` //a310
	Num         int64  `json:"num"`
	AvlNum      int64  `json:"avlNum"`
}

type NodeDeviceTypeOrg struct {
	OrgId int64 `json:"orgId"`
	NodeDeviceType
}

func (d *NodeDevice) genNodeDeviceKey() string {
	switch strings.ToLower(d.Type) {
	case CPU:
		return d.Arch
	default:
		if d.Series == "" {
			return d.Model
		}
		if d.Model == "" {
			return d.Series
		}
		return d.Model + "-" + d.Series
	}
}

func (d *NodeDevice) IsCPU() bool {
	if d.Type == strings.ToUpper(CPU) {
		return true
	}
	return false
}

func (d *NodeDevice) BeforeCreate(tx *gorm.DB) (err error) {
	d.Key = d.genNodeDeviceKey()
	return
}

func (d *NodeDevice) BeforeUpdate(tx *gorm.DB) (err error) {
	d.Key = d.genNodeDeviceKey()
	return
}
