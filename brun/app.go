package main

import (
	"gitee.com/lwx0416/logAgent/infra"
	"gitee.com/lwx0416/logAgent/infra/base"
)

func init() {
	infra.Register(&base.EtcdStarter{})
	infra.Register(&base.KafkaStarter{})

}
