package dto

type VC struct {
	BaseListDto
	ID     int64   `form:"id" json:"id" uri:"id"`
	Name   string  `form:"name" json:"name"`
	Desc   string  `form:"desc" json:"desc"`
	CpuNum float32 `form:"cpuNum" json:"cpuNum"`
	Mem    int64   `form:"mem" json:"mem"`
	Devs   []Dev   `form:"devs" json:"devs"`
}

type Dev struct {
	ID          int64   `form:"id" json:"id" uri:"id"`
	Type        string  `form:"type" json:"type"`
	ModelNumber string  `form:"modelNumber" json:"modelNumber"`
	Num         float32 `form:"num" json:"num"`
	Series      string  `json:"series"`
}
