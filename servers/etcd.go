package servers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"gitee.com/lwx0416/logAgent/infra/base"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

type LogEntry struct {
	logName string `json:"log_name"`
	LogPath string `json:"log_path"`
	Topic   string `json:"topic"`
}

var (
	GlogMgr *LogEntry
)

func (l *LogEntry) WatchTailJob() (err error) {
	var (
		logEntryList []*LogEntry
		log          *LogEntry
		jobEvent     *JobEvent
	)
	client := base.EtcdClient()
	jobWatch := clientv3.NewWatcher(client)
	logAgnetKey, err := GetLocalIP()
	if err != nil {
		logrus.Error(err)
		panic("获取本机ip失败")
	}
	jobKey := JOB_SAVE_DIR + logAgnetKey + "/"
	getResp, err := client.Get(context.TODO(), jobKey, clientv3.WithPrefix())
	if err != nil {
		logrus.Error("从ETCD获取logpath 失败： ", err)
		return
	}
	for _, v := range getResp.Kvs {
		var logEntry *LogEntry
		err = json.Unmarshal(v.Value, &logEntry)
		if err != nil {
			logrus.Error("获取log path 值失败：", err)
			return
		}
		logEntry.logName = string(v.Key)
		logEntryList = append(logEntryList, logEntry)
	}
	logrus.Debug(logEntryList)
	// 初始化任务一次
	for _, v := range logEntryList {
		jobEvent := BuildJobEvent(JOB_EVENT_SAVE, v)
		GScheduler.PushJobEvent(jobEvent)
	}

	// 从当前revision 版本向后监听
	go func() {
		watchStartRevision := getResp.Header.Revision
		watchChan := jobWatch.Watch(context.TODO(), jobKey, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		// 处理监听事件
		for watchResp := range watchChan {
			for _, watchEvent := range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					// 1. 分为添加和修改
					if log, err = UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// 2. 构建event
					jobEvent = BuildJobEvent(JOB_EVENT_SAVE, log)
				case mvccpb.DELETE:
					// 删除
					logName := ExtractJobName(string(watchEvent.Kv.Key))
					log := &LogEntry{logName: logName}
					jobEvent = BuildJobEvent(JOB_EVENT_DELETE, log)
				}
				logrus.Infof("监听事件变化: \t %d , %+v, %+v", jobEvent.EventType, jobEvent.Job)
				fmt.Println("key： ", watchEvent.Kv.Key)
				GScheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}

// 获取本机网卡IP
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)

	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	fmt.Println(addrs)
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4,ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	err = errors.New("没有找到网卡IP")
	return
}
