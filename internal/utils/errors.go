package utils

type CodeMessage struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func ModuleError(baseCode int) int {
	return Module*100000 + baseCode
}

var (
	ActionSuccess        = &CodeMessage{0, "ok.ActionSuccess"}
	ValidateErrorMessage = map[string]string{}
	ErrorNotFound        = &CodeMessage{ModuleError(00001), "err.ResourceNotFound"}
	ErrorValidation      = &CodeMessage{ModuleError(10001), "err.ErrorValidation"}

	ErrorDatabaseOp     = &CodeMessage{ModuleError(20002), "err.ErrorDatabaseOp"}
	ErrorRecordNotExist = &CodeMessage{ModuleError(20005), "err.ErrorRecordNotExist"}
)
