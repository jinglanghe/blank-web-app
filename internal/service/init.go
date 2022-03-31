package service

import (
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/cache"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
)

var (
	nodeDeviceDao        dao.NodeDevices
	resourceQuotaDao     dao.ResourceQuota
)

func Init() {
	if err := initMQ(); err != nil {
		logging.Fatal().Err(err).Msg("init mq failed")
	}

	err := initClientSet()
	if err != nil {
		logging.Fatal().Err(err).Msg("init clientSet failed")
	}

	initInformer()
}

func RefreshData() {
	lockKey := "internal/service/refresh_data"
	mu := cache.RedLock.NewMutex(lockKey)
	err := mu.Lock()
	if err != nil {
		logging.Error(err).Msg("RefreshData failed at lock")
		return
	}
	defer mu.Unlock()

	RefreshModelResourceQuota()
}

func RefreshModelResourceQuota() {
	_types, err := nodeDeviceDao.List()
	if err != nil {
		logging.Fatal().Err(err).Msg("RefreshModelResourceQuota failed at get nodeDeviceType")
		return
	}

	err = resourceQuotaDao.Refresh(_types)
	if err != nil {
		logging.Fatal().Err(err).Msg("RefreshModelResourceQuota failed at get refresh data")
		return
	}
}
