package aaa

import (
	"github.com/apulis/bmod/aistudio-aom/internal/model"
	"time"
)

type User struct {
	model.Base
	Username       string       `gorm:"index:idx_username,unique" json:"username"`
	Password       string       `json:"-"`
	Salt           string       `json:"-"`
	OrganizationID int64        `json:"organizationId"`
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"  json:"organization"`
	UserGroupID    int64        `json:"userGroupId"`
	UserGroup      UserGroup    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"userGroup"`
	CreatorName    string       `gorm:"creator_name" json:"creatorName"`
	FailCount      int          `json:"-"`
	UnlockTime     time.Time    `json:"-"`
	IsInitPassword *bool        `gorm:"is_init_password" json:"-"`
}
