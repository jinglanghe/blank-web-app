package aaa

import "gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"

type Organization struct {
	model.Base
	Account     string `form:"account" json:"account"`
	NickName    string `form:"nickName" json:"nickName"`
	Desc        string `gorm:"desc" json:"desc"`
	VcID        int64  `gorm:"vc_id" json:"vcId"`
	CreatorName string `json:"creatorName"`
	Apps        string
}
