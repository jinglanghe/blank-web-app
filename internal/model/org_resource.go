package model

type OrgResource struct {
	Base
	OrgID       int64  `json:"orgId" gorm:"unique"`
	OrgName     string `json:"orgName"`
	Duration    int64  `json:"duration"`
	Source      int64  `json:"source"`
	CreatorID   int64  `json:"creatorId"`
	CreatorName string `json:"creatorName"`
	Nodes       []Node `json:"nodes"`
}
