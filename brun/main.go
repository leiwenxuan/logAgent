package main

import (
	"gitee.com/lwx0416/logAgent/infra"
	"gitee.com/lwx0416/logAgent/infra/base"
	"github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {

	file := kvs.GetCurrentFilePath("conf.ini", 1)
	logrus.Info("配置文件： ", file)
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
}
