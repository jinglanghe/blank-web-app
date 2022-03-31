package model

type VC struct {
	Base
	Name        string   `gorm:"unique" json:"name"`
	Desc        *string  `json:"desc"`
	CpuNum      *float32 `json:"cpuNum"`
	Mem         *int64   `json:"mem"`
	CreatorName string   `json:"creatorName"`
	Devs        []VCDev  `json:"devs"`
}

type VCDev struct {
	Base
	Type        string  `json:"type"` //
	ModelNumber string  `json:"modelNumber"`
	Series      string  `json:"series"`
	Key         string  `form:"-" json:"-"` // Type + SeriesType作为唯一key
	Num         float32 `json:"num"`
	VCID        int64   `json:"vcId"` // 想设置gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" 但是并没有生效
}

type TotalVCInfo struct {
	CpuNum float32
	Mem    int64
}

type VCStatus struct {
	OrgCount int    `json:"orgCount"`
	OrgName  string `json:"orgName"`
}
