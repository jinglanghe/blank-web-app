package aaa

import (
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gorm.io/datatypes"
)

const (
	PolicyTypeSystemAdmin PolicyType = iota + 1 // 系统管理员组
	PolicyTypeOrgAdmin                          // 组织管理员组
	PolicyTypeDeveloper
	PolicyTypeAnnotator
	PolicyTypeCommon // 普通组织下面的用户组
)

type PolicyType int

type Policy struct {
	model.Base
	Account      string         `gorm:"uniqueIndex" json:"account"`
	Nickname     string         `json:"nickname"`
	Version      string         `json:"version"`
	Desc         string         `json:"desc"`
	Type         string         `json:"type"`
	DbMembers    datatypes.JSON `gorm:"column:members" json:"-"`
	Members      []string       `gorm:"-" json:"-"`
	DbStatements datatypes.JSON `gorm:"column:statements" json:"-"`
	//Statements   []Statement    `gorm:"-" json:"statements"`
	CreatorID   int64      `gorm:"column:creator_id" json:"creatorId"`
	CreatorName string     `json:"creatorName"`
	PolicyType  PolicyType `gorm:"-" json:"policyType"`
}

type InitPolicies struct {
	Module    string         `json:"module"`
	SysAdmin  []StatementDto `json:"systemAdmin"`
	OrgAdmin  []StatementDto `json:"orgAdmin"`
	Developer []StatementDto `json:"developer"`
}

type StatementDto struct {
	Actions   []string `form:"actions"`
	Resources []string `form:"resources"`
	Effect    string   `form:"effect" validate:"oneof=allow deny"`
	Role      string   `form:"role"`
}
