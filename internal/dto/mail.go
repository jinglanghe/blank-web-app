package dto

type MailSendDto struct {
	OrgId    *int64 `json:"orgId"`
	Host     string `json:"host"`
	Port     int    `json:"port" binding:"required_with=Host"`
	UserName string `json:"userName" binding:"required_with=Host"`
	Password string `json:"password" binding:"required_with=Host"`

	From    string `json:"from"`
	To      string `json:"to" binding:"required"`
	Cc      string `json:"cc"`
	Bcc     string `json:"bcc"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}
