package service

import (
	"encoding/json"
	"time"

	config "github.com/apulis/bmod/aistudio-aom/configs"
	"github.com/apulis/bmod/aistudio-aom/internal/cache"
	"github.com/apulis/bmod/aistudio-aom/internal/dao"
	"github.com/apulis/bmod/aistudio-aom/internal/metadata"
	"github.com/apulis/bmod/aistudio-aom/internal/model/aaa"
	aaa2 "github.com/apulis/bmod/aistudio-aom/internal/service/aaa"
	"github.com/apulis/sdk/go-utils/logging"
)

var (
	nodeDao              dao.Nodes
	nodeDeviceDao        dao.NodeDevices
	userGroupResourceDao dao.UserGroupResource
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

	//go RegisterEndpointsAndPolicies()
	//go IntervalRefresh()

	initInformer()
}

func RegisterEndpointsAndPolicies() {
	endpointBytes, err := metadata.Asset("endpoint.json")
	if err != nil {
		logging.Fatal().Err(err).Msg("load endpoint json file error")
		return
	}

	var endpoints aaa.EndPointsAndPolicies
	if err := json.Unmarshal(endpointBytes, &endpoints); err != nil {
		logging.Fatal().Err(err).Msg("unmarshal endpoint json file error")
		return
	}

	if err := aaa2.RegisterEndPoints(&endpoints); err != nil {
		logging.Error(err).Msg("RegisterEndPoints error")
	}
}

func IntervalRefresh() {
	nodeInfoInterval := config.Config.NodeInfoInterval
	if nodeInfoInterval == 0 {
		nodeInfoInterval = 60
	}
	logging.Info().Msgf("refresh interval: %v minute", nodeInfoInterval)
	for {
		time.Sleep(time.Duration(nodeInfoInterval) * time.Minute)
		RefreshData()
	}
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

	RefreshModelNode()
	RefreshNodeLabel()
	RefreshQuota()
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
