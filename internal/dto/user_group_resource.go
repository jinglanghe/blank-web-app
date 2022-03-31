package dto

type UserGroupResourceCreate struct {
	OrgID         int64  `json:"orgId" binding:"required"`
	OrgName       string `json:"orgName"`
	UserGroupID   int64  `json:"userGroupId" binding:"required"`
	UserGroupName string `json:"userGroupName"`
	Min           *int64 `json:"min" binding:"required,min=0,max=100,ltefield=Max"`
	Max           *int64 `json:"max" binding:"required,min=0,max=100"`
	Duration      int64  `json:"duration"`
	Source        int64  `json:"source"`
}

type UserGroupResourceList struct {
	BaseListDto
	UserGroupId int64 `uri:"userGroupId" form:"userGroupId" json:"userGroupId"`
}

type UserGroupResourceGet struct {
	ID int64 `uri:"id" json:"id" binding:"required"`
}

type UserGroupResourceUpdate struct {
	ID  int64  `uri:"id" json:"id" binding:"required"`
	Min *int64 `json:"min" binding:"required,min=0,max=100,ltefield=Max"`
	Max *int64 `json:"max" binding:"required,min=0,max=100"`
}

type UserGroupResourceDelete struct {
	ID int64 `uri:"id" json:"id" binding:"required"`
}

type UserGroupResourcePreview struct {
	ID  int64  `uri:"id" json:"id" binding:"required"`
	Min *int64 `uri:"min" form:"min" json:"min" binding:"required,min=0,max=100,ltefield=Max"`
	Max *int64 `uri:"max" form:"max" json:"max" binding:"required,min=0,max=100"`
}
