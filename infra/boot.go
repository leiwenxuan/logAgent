package infra

import (
	"github.com/tietang/props/kvs"
)
// 应用程序
type BootApplication struct {
	IsTest bool
	conf kvs.ConfigSource
	starterCtx StarterContext
}

// 构造系统, 初始化配置文件
func New(conf kvs.ConfigSource) *BootApplication  {
	e := &BootApplication{conf:conf, starterCtx:StarterContext{}}
	e.starterCtx.SetProps(conf)
	return e
}
func (b *BootApplication) Start()  {
	// 1. 初始化
	b.init()
	// 2. 安装
	b.setup()
	// 3. 启动
	b.start()

}

func (e *BootApplication)init()  {
	for _, v := range GetStarters() {
		v.Init(e.starterCtx)
	}
}

func (e *BootApplication)setup()  {
	for _, v := range GetStarters() {
		v.Setup(e.starterCtx)
	}
}

func (e *BootApplication)start()  {
	for i, v := range GetStarters() {
		if v.StartBlocking(){
			if i+1 == len(GetStarters()){
				v.Start(e.starterCtx)
			}else {
				go v.Start(e.starterCtx)
			}
		}else {
			v.Start(e.starterCtx)
		}
	}
}

func (e *BootApplication) Stop()  {
	for _, v:= range GetStarters(){
		v.Stop(e.starterCtx)
	}
}



