package servers

import (
	"context"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

// 日志收集模块
type TailJobMgr struct {
	instance   *tail.Tail
	path       string
	topic      string
	ctx        context.Context
	cancelFunc context.CancelFunc
}

var (
	GTailJobServer *TailJobMgr
)

func NewTailJob(path string, topic string) (tailJob *TailJobMgr) {
	ctx, cancel := context.WithCancel(context.TODO())
	tailJob = &TailJobMgr{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	if err := tailJob.Init(); err != nil {
		logrus.Error("打开文件失败： ", err)
	}
	return
}

// 初始化 tail 客户端, 返回单利对象
func (t *TailJobMgr) Init() (err error) {
	config := tail.Config{
		ReOpen:    true,  // 是否重新打开
		MustExist: false, // 文件不存在不报错， 防止log文件没有生成是报错
		Poll:      true,  //
		Follow:    true,  // 是否跟随，
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	}
	logrus.Debug("打开的日志文件： ", t.path)
	//filename := "F:\\code\\002Golang\\go_demo\\tail-demo\\tail-01\\my.log"
	//t.path = filename
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		logrus.Error("打开log 文件失败： ", err)
		return
	}
	// 发送日志
	go t.run()
	return
}

func (t *TailJobMgr) ReadChan() <-chan *tail.Line {
	return t.instance.Lines
}

func (t *TailJobMgr) run() {
	var line *tail.Line
	for {
		select {
		case <-t.ctx.Done():
			// 配置更新退出程序
			return
		case line = <-t.instance.Lines:
			logrus.Debug(line.Text)
			SendToChan(t.topic, line.Text)
		}
	}

}
