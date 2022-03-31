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
	ErrorMissingOrgId    = &CodeMessage{ModuleError(00002), "err.ErrorMissingOrgId"}
	ErrorMissingJwtToken = &CodeMessage{ModuleError(00003), "err.ErrorMissingJwtToken"}
	ErrorValidation      = &CodeMessage{ModuleError(10001), "err.ErrorValidation"}

	ErrUserCurrent               = &CodeMessage{ModuleError(10701), "ErrUserCurrent"}
	ErrOrgDetail                 = &CodeMessage{ModuleError(10711), "ErrOrgDetail"}
	ErrUserGroupDetail           = &CodeMessage{ModuleError(10721), "ErrUserGroupDetail"}
	ErrUserGroupInvalid          = &CodeMessage{ModuleError(10722), "ErrUserGroupInvalid"}
	ErrNodeDeviceTypeNotExist    = &CodeMessage{ModuleError(10731), "ErrNodeDeviceTypeNotExist"}
	ErrResourceQuotaAlreadyExist = &CodeMessage{ModuleError(10741), "ErrResourceQuotaAlreadyExist"}

	ErrUserGroupResourceExceed     = &CodeMessage{ModuleError(10801), "ErrUserGroupResourceExceed"}
	ErrUserGroupResourceMinInvalid = &CodeMessage{ModuleError(10802), "ErrUserGroupResourceMinInvalid"}

	ErrorDatabaseOp     = &CodeMessage{ModuleError(20002), "err.ErrorDatabaseOp"}
	ErrorRuleValidate   = &CodeMessage{ModuleError(20003), "err.ErrorRuleValidate"}
	ErrorMailFailed     = &CodeMessage{ModuleError(20004), "err.ErrorMailFailed"}
	ErrorRecordNotExist = &CodeMessage{ModuleError(20005), "err.ErrorRecordNotExist"}
	ErrorSettingGet     = &CodeMessage{ModuleError(20006), "err.ErrorSettingGet"}
	ErrorAesDecrypt     = &CodeMessage{ModuleError(20007), "err.ErrorAesDecrypt"}
	ErrorAesEncrypt     = &CodeMessage{ModuleError(20008), "err.ErrorAesEncrypt"}
	ErrorConfigmapOp    = &CodeMessage{ModuleError(20009), "err.ErrorConfigmapOp"}
	ErrorConfigmap404   = &CodeMessage{ModuleError(20010), "err.ErrorConfigmapOp"}

	ErrModelArtsCreate = &CodeMessage{ModuleError(40001), "ErrModelArtsCreate"}
	ErrModelArtsList   = &CodeMessage{ModuleError(40002), "ErrModelArtsList"}
	ErrZipEntryError   = &CodeMessage{ModuleError(40003), "ErrZipEntryError"}
	ErrZipWriteError   = &CodeMessage{ModuleError(40004), "ErrZipWriteError"}

	ErrNodeCordon   = &CodeMessage{ModuleError(60001), "ErrNodeCordon"}
	ErrNodeUnCordon = &CodeMessage{ModuleError(60002), "ErrNodeUnCordon"}
	ErrNodePodsList = &CodeMessage{ModuleError(60011), "ErrNodePodsList"}
	ErrNodeRequests = &CodeMessage{ModuleError(60021), "ErrNodeRequests"}
	ErrNodeLabelSet = &CodeMessage{ModuleError(60031), "ErrNodeLabelSet"}
	ErrNodeTaintAdd = &CodeMessage{ModuleError(60041), "ErrNodeTaintAdd"}

	ErrNamespaceUsedRes = &CodeMessage{ModuleError(61001), "ErrNamespaceUsedRes"}

	ErrMailSend = &CodeMessage{ModuleError(70001), "ErrMailSend"}
)
