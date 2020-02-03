package base

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"gitee.com/lwx0416/logAgent/infra"
	"github.com/Shopify/sarama"
)

var clientKafka sarama.SyncProducer

type KafkaStarter struct {
	infra.BaseStarter
}

func KakfaClient() sarama.SyncProducer {
	Check(clientEtcd)
	return clientKafka
}

func (k *KafkaStarter) Setup(ctx infra.StarterContext) {
	var ()
	conf := ctx.Props()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	addrs := conf.GetDefault("kafka.endpoints", "127.0.0.1:9092")
	if addrs == "" {
		panic("获取kafka集群失败")
	}
	addrsList := strings.Split(addrs, ",")
	logrus.Info("kafka集群: ", addrsList)
	client, err := sarama.NewSyncProducer(addrsList, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	clientKafka = client
}
