package model

type UserClaims struct {
	UserID         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	GroupID        int64  `json:"group_id"`
	GroupAccount   string `json:"group_account"`
	OrganizationId int64  `json:"organization_id"`
}
