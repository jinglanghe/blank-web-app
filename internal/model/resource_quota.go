package model

type ResourceQuota struct {
	Base
	Type        string `json:"type"`
	Key         string `json:"key"`
	Arch        string `json:"arch"`
	Model       string `json:"model"`
	ComputeType string `json:"computeType"`
	Series      string `json:"series"`
	CpuNum      int64  `json:"cpuNum"`
	Mem         int64  `json:"mem"`
	CreatorID   int64  `json:"creatorId"`
	CreatorName string `json:"creatorName"`
	Num         int64  `json:"num"`
	AvlNum      int64  `json:"avlNum"`
}

func (r *ResourceQuota) EqualNodeDevice(device NodeDevice) bool {
	if r.Type == device.Type &&
		r.Key == device.Key &&
		r.Arch == device.Arch &&
		r.Model == device.Model &&
		r.ComputeType == device.ComputeType &&
		r.Series == device.Series {
		return true
	}
	return false
}

func (r *ResourceQuota) EqualNodeDeviceType(deviceType *NodeDeviceType) bool {
	if r.Type == deviceType.Type &&
		r.Key == deviceType.Key &&
		r.Arch == deviceType.Arch &&
		r.Model == deviceType.Model &&
		r.ComputeType == deviceType.ComputeType &&
		r.Series == deviceType.Series {
		return true
	}
	return false
}
