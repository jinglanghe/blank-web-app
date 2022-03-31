package model

type ModelArts struct {
	Base
	UserName   string `json:"userName"`
	AK         string `json:"aK"`
	SK         string `json:"sk"` // 这里失误了， 其实本应该想写成sK的
	ProjectId  string `json:"projectId"`
	BucketName string `json:"bucketName"`
	Site       string `json:"site"`
}
