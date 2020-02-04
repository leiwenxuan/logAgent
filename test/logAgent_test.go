package test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"

	"gitee.com/lwx0416/logAgent/infra"
	"gitee.com/lwx0416/logAgent/infra/base"
	"gitee.com/lwx0416/logAgent/servers"
)

func TestGetEtcdConfPath(t *testing.T) {
	infra.Register(&base.EtcdStarter{})
	infra.Register(&base.KafkaProducerStarter{})

	file := kvs.GetCurrentFilePath("conf.ini", 1)
	logrus.Info("配置文件： ", file)
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
	_ = servers.GlogMgr.WatchTailJob()
}
