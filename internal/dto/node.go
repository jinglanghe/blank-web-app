package dto

type NodeList struct {
	BaseListDto
	Type            string `json:"type" form:"type" uri:"type"`
	CpuArch         string `json:"cpuArch" form:"cpuArch" uri:"cpuArch"`
	CpuNum          int64  `json:"cpuNum" form:"cpuNum" uri:"cpuNum"`
	AvlCpuNum       int64  `json:"avlCpuNum" form:"avlCpuNum" uri:"avlCpuNum"`
	Mem             int64  `json:"mem" form:"mem" uri:"mem"`
	AvlMem          int64  `json:"avlMem" form:"avlMem" uri:"avlMem"`
	Role            string `json:"role" form:"role" uri:"role"`
	Status          string `json:"status" form:"status" uri:"status"`
	Keyword         string `json:"keyword" form:"keyword" uri:"keyword"`
	ResourceQuotaID int64  `json:"resourceQuotaId" form:"resourceQuotaId" uri:"resourceQuotaId"`
}

type NodePodList struct {
	Ids string `json:"ids" uri:"ids" form:"ids"`
}

type NodeLabelCreate struct {
	ID    int64  `uri:"id" form:"id" json:"id" binding:"required"`
	Key   string `uri:"key" form:"key" json:"key" binding:"required"`
	Value string `uri:"value" form:"value" json:"value" binding:"required"`
}

type NodeTaintCreate struct {
	ID     int64  `uri:"id" form:"id" json:"id" binding:"required"`
	Key    string `uri:"key" form:"key" json:"key" binding:"required"`
	Value  string `uri:"value" form:"value" json:"value" binding:"required"`
	Effect string `uri:"effect" form:"effect" json:"effect" binding:"required"`
}
