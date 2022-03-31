package model

type Setting struct {
	Base
	OrgId          int64  `gorm:"column:org_id" json:"-"`
	SenderMail     string `gorm:"column:sender_mail" json:"email"`
	SenderPassword string `gorm:"column:sender_password" json:"password"`
	SmtpServer     string `gorm:"smtp_server" json:"smtpServer"`
	MailType       string `json:"type"`
	Receiver       string `json:"to"`
}

func (Setting) TableName() string {
	return "alert_setting"
}
