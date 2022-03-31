package service

import (
	"fmt"
	"github.com/apulis/bmod/aistudio-aom/configs"
	"github.com/apulis/sdk/go-utils/broker"
	_ "github.com/apulis/sdk/go-utils/broker/kafka"
	"github.com/apulis/sdk/go-utils/broker/rabbitmq"
	"github.com/apulis/sdk/go-utils/logging"
)

var (
	addr  = "amqp://guest:guest@localhost:5672/"
	topic = ""
	mq    broker.Broker
)

const (
	TYPE = "type"
	NS   = "ns"
)

func initMQ() error {
	if mq != nil {
		return nil
	}

	topic = configs.Config.Rabbitmq.Topic
	addr = "amqp://" + configs.Config.Rabbitmq.Username + ":" + configs.Config.Rabbitmq.Password + "@" + configs.Config.Rabbitmq.Host + ":" + fmt.Sprintf("%d", configs.Config.Rabbitmq.Port) + "/"
	mq = rabbitmq.NewBroker(
		broker.Addrs(addr),
		rabbitmq.ExchangeName("k8s_exchange"),
	)

	if err := mq.Connect(); err != nil {
		return err
	}

	f := func(event broker.Event) error {
		logging.Debug().Msg("receive mq")
		if err := handler(event.Message().Header); err != nil {
			logging.Error(err).Msg("process msg error")
			return err
		}
		return nil
	}

	go func() {
		_, err := mq.Subscribe(topic, f, rabbitmq.DurableQueue())
		if err != nil {
			logging.Error(err).Msg("receive mq failed")
		}
	}()

	return nil
}

func handler(m map[string]string) error {
	typeStr := ""
	ns := ""

	for k, v := range m {
		if k == TYPE {
			typeStr = v
		} else if k == NS {
			ns = v
		}
	}

	if typeStr != "1" || len(ns) == 0 {
		logging.Warn().Msgf("unknown mq header: %+v\n", m)
		return nil
	}

	// 创建ns
	CreateNamespace(ns)

	return nil
}
