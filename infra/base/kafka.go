package base

import (
	"strings"

	"github.com/sirupsen/logrus"

	"gitee.com/lwx0416/logAgent/infra"
	"github.com/Shopify/sarama"
)

var kafkaProducer sarama.SyncProducer
var kafkaConsumer sarama.Consumer

type KafkaProducerStarter struct {
	infra.BaseStarter
}
type KafkaConsumerStarter struct {
	infra.BaseStarter
}

func KakfaProducerClient() sarama.SyncProducer {
	Check(kafkaProducer)
	return kafkaProducer
}

func KakfaConsumerClient() sarama.Consumer {
	Check(kafkaConsumer)
	return kafkaConsumer
}

func (k *KafkaProducerStarter) Setup(ctx infra.StarterContext) {
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
		panic("producer close, err:" + err.Error())
		return
	}
	kafkaProducer = client
}

func (k *KafkaConsumerStarter) Setup(ctx infra.StarterContext) {
	var ()
	conf := ctx.Props()

	addrs := conf.GetDefault("kafka.endpoints", "127.0.0.1:9092")
	if addrs == "" {
		panic("获取kafka集群失败")
	}
	addrsList := strings.Split(addrs, ",")
	logrus.Info("kafka集群: ", addrsList)
	client, err := sarama.NewConsumer(addrsList, nil)
	if err != nil {
		panic("producer close, err:" + err.Error())
		return
	}
	kafkaConsumer = client
}
