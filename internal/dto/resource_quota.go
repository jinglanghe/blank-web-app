package dto

type ResourceQuotaCreate struct {
	Type   string `json:"type" binding:"required"`
	Key    string `json:"key" binding:"required"`
	CpuNum int64  `json:"cpuNum" binding:"required,min=0"`
	Mem    int64  `json:"mem" binding:"required,min=0"`
}

type ResourceQuotaList struct {
	BaseListDto
	Type        string `json:"type" form:"type" uri:"type"`
	CpuNum      int64  `json:"cpuNum" form:"cpuNum" uri:"cpuNum"`
	Mem         int64  `json:"mem" form:"mem" uri:"mem"`
	Arch        string `json:"arch" form:"arch" uri:"arch"`
	AvlNum      int64  `json:"avlNum" form:"avlNum" uri:"avlNum"`
	CreatorName string `json:"creatorName" form:"creatorName" uri:"creatorName"`
}

type ResourceQuotaGet struct {
	ID int64 `uri:"id" json:"id" binding:"required"`
}

type ResourceQuotaDelete struct {
	ID int64 `uri:"id" json:"id" binding:"required"`
}
