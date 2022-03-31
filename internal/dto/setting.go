package dto

type SenderSettingDto struct {
	OrgId      int64  `json:"-"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	SmtpServer string `json:"smtpServer"`
	MailType   string `json:"type"`
	Receiver   string `json:"to"`
}

type OrgSendMailDto struct {
	Email string `form:"email" json:"email"`
}

type SettingTestMailDto struct {
	From       string `json:"from"`
	Password   string `json:"password"`
	SmtpServer string `json:"smtpServer"`
	To         string `json:"to"`
}
