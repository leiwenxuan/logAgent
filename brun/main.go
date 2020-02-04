package main

import (
	"gitee.com/lwx0416/logAgent/infra"
	"gitee.com/lwx0416/logAgent/infra/base"
	"gitee.com/lwx0416/logAgent/servers"
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

	// 初始化监听器
	_ = servers.GlogMgr.WatchTailJob()
	// 启动kafka, 初始化
	go servers.InitKafkaProducer()
	// TODO 创建Kafka分区
	servers.InitEs(5)

	// 初始资源阻塞进程
	ch := make(chan int, 1)
	<-ch

}
