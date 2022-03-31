package dto

type Condition struct {
	Op    string `json:"operator"`
	Value int    `json:"value"`
}

type ResourceDataFilter struct {
	Range int    `json:"range"`
	Type  string `json:"type"`
}

type ServiceDataFilter struct {
	Range      int    `json:"range"`
	Process    string `json:"process"`
	SubProcess string `json:"subProcess"`
}

type ResourceRuleDto struct {
	Name         string             `json:"name" binding:"required"`
	Source       string             `json:"source" binding:"required"`
	Type         string             `json:"type" binding:"required"`
	DataFilter   ResourceDataFilter `json:"dataFilter" binding:"required"`
	Condition    Condition          `json:"condition" binding:"required"`
	Level        string             `json:"level" binding:"required"`
	Frequency    int                `json:"frequency" binding:"required"`
	TriggerCount int                `json:"triggerCount" binding:"required"`
	Interval     int                `json:"interval"`
	Title        string             `json:"title" binding:"required"`
	Content      string             `json:"content" binding:"required"`
	Receivers    string             `json:"receivers"`
	Status       int                `json:"status"`
}

type ResourceRuleCreateDtd struct {
	ResourceRuleDto
}

type ResourceRuleEditDto struct {
	ResourceRuleDto
	UUID string `uri:"id" json:"id" binding:"required"`
}

type ServiceRuleDto struct {
	Name         string            `json:"name" binding:"required"`
	Source       string            `json:"source" binding:"required"`
	Type         string            `json:"type" binding:"required"`
	DataFilter   ServiceDataFilter `json:"dataFilter" binding:"required"`
	Condition    Condition         `json:"condition" binding:"required"`
	Level        string            `json:"level" binding:"required"`
	Frequency    int               `json:"frequency" binding:"required"`
	TriggerCount int               `json:"triggerCount" binding:"required"`
	Interval     int               `json:"interval"`
	Title        string            `json:"title" binding:"required"`
	Content      string            `json:"content" binding:"required"`
	Receivers    string            `json:"receivers"`
	Status       int               `json:"status"`
}

type ServiceRuleCreateDtd struct {
	ServiceRuleDto
}

type ServiceRuleEditDto struct {
	ServiceRuleDto
	UUID string `uri:"id" json:"id" binding:"required"`
}

type RuleStatusDto struct {
	ID     string `uri:"id" json:"id" binding:"required"`
	Status int    `json:"status"`
}

type RuleGetDto struct {
	ID string `uri:"id" json:"id" binding:"required"`
}

type RuleDeleteDto struct {
	ID string `uri:"id" json:"id" binding:"required"`
}

type RuleListDto struct {
	BaseListDto
	Name      string `form:"name" json:"name"`
	StartTime int64  `form:"startTime" json:"startTime" time_format:"unix"`
	EndTime   int64  `form:"endTime" json:"endTime" time_format:"unix"`
}
