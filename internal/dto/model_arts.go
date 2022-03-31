package dto

type ModelArts struct {
	UserName   string `form:"userName" binding:"required"`
	AK         string `form:"aK" binding:"required"`
	SK         string `form:"sk" binding:"required"`
	ProjectId  string `form:"projectId" binding:"required"`
	BucketName string `form:"bucketName" binding:"required"`
	Site       string `form:"site" binding:"required"`
}
