package model

const (
	CPU         = "cpu"
	GPU         = "gpu"
	NPU         = "npu"
	SERIES      = "series"
	ComputeType = "computeType"
)

const (
	NodeStatusReady    = "Ready"
	NodeStatusNotReady = "NotReady"
)

type Node struct {
	Base
	OrgResourceId int64        `json:"orgResourceId"`
	MachineID     string       `json:"machineId" gorm:"unique"`
	Name          string       `json:"name"`
	InternalIP    string       `json:"internalIp"`
	Type          string       `json:"type"`    //CPU|GPU|NPU
	CpuArch       string       `json:"cpuArch"` //amd64|x86
	CpuNum        float32      `json:"cpuNum"`
	AvlCpuNum     float32      `json:"avlCpuNum"`
	Mem           int64        `json:"mem"`
	AvlMem        int64        `json:"avlMem"`
	Role          string       `json:"role"`
	Status        string       `json:"status"`
	Devs          []NodeDevice `json:"devs" `
	OrgResource   OrgResource  `json:"orgResource"`
}

type OrgAvlMem struct {
	OrgId  int64 `json:"orgId"`
	AvlMem int64 `json:"avlMem"`
}

type NodeOverview struct {
	ResStats ResStats `json:"resStats"`
}

type ResStats struct {
	Cpu NodeTuple `json:"cpu"`
	Gpu NodeTuple `json:"gpu"`
	Npu NodeTuple `json:"npu"`
}

type NodeTuple struct {
	Occupied  int64 `json:"occupied"`
	Available int64 `json:"available"`
}

type PodInfo struct {
	Id     int64
	PodNum int `json:"podNum"`
}
