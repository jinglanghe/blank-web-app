package aaa

type EndPointsAndPolicies struct {
	EndPoints []EndPoint   `json:"endpoints"`
	Policies  InitPolicies `json:"policies"`
}

type EndPoint struct {
	Module       string `form:"module"`
	Desc         string `form:"desc"`
	Resource     string `form:"resource"`
	Action       string `form:"action"`
	HttpMethod   string `form:"httpMethod"`
	HttpEndpoint string `form:"httpEndpoint"`
}
