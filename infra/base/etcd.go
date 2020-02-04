package base

import (
	"strings"
	"time"

	"gitee.com/lwx0416/logAgent/infra"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

//dbx 数据库实例
var clientEtcd *clientv3.Client

func EtcdClient() *clientv3.Client {
	Check(clientEtcd)
	return clientEtcd
}

//Etcd starter，并且设置为全局
type EtcdStarter struct {
	infra.BaseStarter
}

func (s *EtcdStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	endpoints := conf.GetDefault("etcd.endpoints", "")
	if endpoints == "" {
		panic("获取ETCD集群失败")
	}
	endpointsList := strings.Split(endpoints, ",")
	log.Info("获取ETCD集群：", endpointsList)
	etcdDialTimeout := conf.GetIntDefault("etcd.dialTimeout", 5000)

	log.Info("获取ETCD超时时间：", etcdDialTimeout)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   endpointsList,                                     // 集群地址
		DialTimeout: time.Duration(etcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		panic("连接失败")
		return
	}

	clientEtcd = client
	log.Debug("client etcd:", clientEtcd)
	return
}
