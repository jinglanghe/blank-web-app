package aaa

import "github.com/apulis/bmod/aistudio-aom/internal/model"

const (
	SystemAdmin UserGroupType = iota + 1 // 系统管理员组
	OrgAdmin                             // 组织管理员组
	Developer
	Annotator
	Common   // 普通组织下面的用户组
	OrgEmpty // 没有分配到组织
)

type UserGroupType int

type UserGroup struct {
	model.Base
	Account         string        `gorm:"unique" json:"account"`
	OrganizationID  int64         `gorm:"unique" json:"organizationId"`
	Organization    Organization  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"organization"`
	NickName        string        `json:"nickName"`
	Desc            string        `gorm:"desc" json:"desc"`
	ResourceQuotaId int64         `json:"resourceQuotaId"`
	Policies        []Policy      `gorm:"many2many:ug_policies;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"policies"`
	CreatorID       int64         `gorm:"creator_id" json:"creatorId"`
	CreatorName     string        `json:"creatorName"`
	Type            UserGroupType `gorm:"-" json:"type"`
}
