package service

import (
	_ "gitlab.apulis.com.cn/hjl/blank-web-app/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app/logging"
)

//var (
//	nodeDeviceDao        dao.NodeDevices
//	resourceQuotaDao     dao.ResourceQuota
//)

func Init() {
	if err := initMQ(); err != nil {
		logging.Fatal().Err(err).Msg("init mq failed")
	}
}
