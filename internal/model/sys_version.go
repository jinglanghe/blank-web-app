package model

type SysVersion struct {
	Base
	SysVersion  string `json:"sysVersion"`
	InstallTime int64  `json:"installTime"`
	Installer   string `json:"installer"`
	Description string `json:"description"`
}
