package dto

type AlertRecordListDto struct {
	BaseListDto
	StartTime int64  `form:"startTime" json:"startTime" time_format:"unix" binding:"required"`
	EndTime   int64  `form:"endTime" json:"endTime" time_format:"unix" binding:"required"`
	Level     int    `form:"level" json:"level"`
	Name      string `form:"name" json:"name"`
}
