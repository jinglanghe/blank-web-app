package aaa

import "github.com/apulis/bmod/aistudio-aom/internal/model"

type Organization struct {
	model.Base
	Account     string `form:"account" json:"account"`
	NickName    string `form:"nickName" json:"nickName"`
	Desc        string `gorm:"desc" json:"desc"`
	VcID        int64  `gorm:"vc_id" json:"vcId"`
	CreatorName string `json:"creatorName"`
	Apps        string
}
