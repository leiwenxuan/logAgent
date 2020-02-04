package main

import (
	"gitee.com/lwx0416/logAgent/infra"
	"gitee.com/lwx0416/logAgent/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.EtcdStarter{})
	infra.Register(&base.KafkaProducerStarter{})
	infra.Register(&base.KafkaConsumerStarter{})
	infra.Register(&base.ESStarter{})

}
