package base

import (
	"gitee.com/lwx0416/logAgent/infra"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
)

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	Check(props)
	return props
}

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	log.Info("初始化配置.")
}
