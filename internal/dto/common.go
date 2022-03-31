package dto

type BaseListDto struct {
	OrgId     int64  `form:"orgId" uri:"orgId" json:"orgId"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
	StartTime int64  `form:"startTime" json:"startTime" time_format:"unix"`
	EndTime   int64  `form:"endTime" json:"endTime" time_format:"unix"`
	Sort      string `form:"sort"`
}
