package dto

type OrgResourceCreate struct {
	OrgName     string  `json:"orgName"`
	OrgId       int64   `uri:"orgId" json:"orgId" binding:"required"`
	NodeIds     []int64 `json:"nodeIds"`
	Duration    int64   `json:"duration"`
	Source      int64   `json:"source"`
	CreatorID   int64   `json:"creatorId"`
	CreatorName string  `json:"creatorName"`
}

type OrgResourceList struct {
	BaseListDto
}

type OrgResourceGet struct {
	ID int64 `uri:"id" json:"id" binding:"required"`
}

type OrgResourceDelete struct {
	ID int64 `uri:"id" json:"id" binding:"required"`
}
