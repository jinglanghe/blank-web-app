package dao

import (
	"encoding/base64"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dto"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/utils"
)

type Setting struct{}

func (s *Setting) SetSender(sender *dto.SenderSettingDto) error {
	password, err := utils.AesEncrypt([]byte(sender.Password), utils.DefaultKey)
	if err != nil {
		return err
	}

	pass64 := base64.StdEncoding.EncodeToString(password)
	update := model.Setting{
		SenderMail:     sender.Email,
		SenderPassword: pass64,
		SmtpServer:     sender.SmtpServer,
		MailType:       sender.MailType,
		Receiver:       sender.Receiver,
	}
	insert := model.Setting{
		OrgId:          sender.OrgId,
		SenderMail:     sender.Email,
		SenderPassword: pass64,
		SmtpServer:     sender.SmtpServer,
		MailType:       sender.MailType,
		Receiver:       sender.Receiver,
	}
	return GetDB().Where("org_id = ?", sender.OrgId).Assign(update).FirstOrCreate(&insert).Error
}

func (s *Setting) GetSender(orgId int64) (*model.Setting, error) {
	var setting model.Setting
	db := GetDB().Where("org_id = ?", orgId).First(&setting)
	if db.Error != nil {
		return nil, db.Error
	}
	bytesPass, _ := base64.StdEncoding.DecodeString(setting.SenderPassword)
	password, err := utils.AesDecrypt(bytesPass, utils.DefaultKey)
	if err != nil {
		return nil, err
	}
	setting.SenderPassword = string(password)
	return &setting, nil
}
