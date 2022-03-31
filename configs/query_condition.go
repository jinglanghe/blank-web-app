package configs

var (
	QueryConditionMetaData QueryConditionMeta
)

type QueryConditionMeta struct {
	Roles      []ValueTuple  `json:"roles"`
	Status     []ValueTuple  `json:"status"`
	Source     []numberTuple `json:"source"`
	Types      interface{}   `json:"types"`
	AlertTypes []ValueTuple  `json:"alertTypes"`
}

type numberTuple struct {
	Value int64  `json:"value"`
	Label string `json:"label"`
}
