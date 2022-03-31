package configs

type MetadataConf struct {
	Resource        resourceMeta          `json:"resourceRule"`
	Service         serviceMeta           `json:"serviceRule"`
	CondOperator    map[string]ValueTuple `json:"conditionOp"`
	Level           map[string]ValueTuple `json:"level"`
	ReceiverChannel map[string]ValueTuple `json:"receiverChannel"`
}

type resourceMeta struct {
	ComputeType map[string]ValueTuple `json:"computeType"`
	Source      map[string]sourceMeta `json:"source"`
}

type serviceMeta struct {
	Source    map[string]sourceMeta `json:"source"`
	DataRange map[string]ValueTuple `json:"dataRange"`
}

type sourceMeta struct {
	Label string                `json:"label"`
	Types map[string]ValueTuple `json:"types"`
}

type ValueTuple struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
