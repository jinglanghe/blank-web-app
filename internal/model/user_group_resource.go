package model

import "strings"

type UserGroupResource struct {
	Base
	OrgID         int64  `json:"orgId"`
	OrgName       string `json:"orgName"`
	UserGroupID   int64  `json:"userGroupId" gorm:"unique"`
	UserGroupName string `json:"userGroupName"`
	Min           int64  `json:"min"`
	Max           int64  `json:"max"`
	Duration      int64  `json:"duration"`
	Source        int64  `json:"source"`
	CreatorID     int64  `json:"creatorId"`
	CreatorName   string `json:"creatorName"`
}

type UserGroupResourceList struct {
	UserGroupResource
	Used   int64   `json:"used"`
	Mem    int64   `json:"mem"`
	MemMin int64   `json:"memMin"`
	MemMax int64   `json:"memMax"`
	Quotas []Quota `json:"quotas"`
}

type Quota struct {
	Type        string `json:"type"` //CPU|GPU|NPU
	Key         string `json:"key"`
	Arch        string `json:"arch"`
	Model       string `json:"model"`
	ComputeType string `json:"computeType"`
	Series      string `json:"series"`
	Num         int64  `json:"num"`
	Used        int64  `json:"used"`
	Min         int64  `json:"min"`
	Max         int64  `json:"max"`
}

func (q *Quota) IsCPU() bool {
	if q.Type == strings.ToUpper(CPU) {
		return true
	}
	return false
}

func (u *UserGroupResourceList) Namespace() string {
	return u.UserGroupName
}

func (u *UserGroupResourceList) QuotaName() string {
	return u.UserGroupName
}
