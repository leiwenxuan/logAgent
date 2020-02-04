package servers

import (
	"context"
	"fmt"

	"gitee.com/lwx0416/logAgent/infra/base"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type logData struct {
	topic string
	msg   string
}

var (
	logDataChan chan *logData
)

func InitKafkaProducer() {
	conf := base.Props()
	chanSize := conf.GetIntDefault("kafka.chanMsgSize", 100000)
	fmt.Println(chanSize)
	logDataChan = make(chan *logData, 100000)
	go SendToKafka()
}

type KafkaServer struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func SendToKafka() {
	for {
		select {
		case msgChan := <-logDataChan:
			msg := &sarama.ProducerMessage{
				Topic: msgChan.topic,
				Value: sarama.StringEncoder(msgChan.msg),
			}
			client := base.KakfaProducerClient()
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Error("Send msg failed, err: ", err)
			}
			logrus.Debug("发送消息成功： ", pid, offset, msgChan.topic)

		default:
		}
	}
}
func SendToChan(topic string, data string) {
	sendMsg := &logData{
		msg:   data,
		topic: topic,
	}
	logDataChan <- sendMsg
}

func TransferToES(tailJob *TailJobMgr) (err error) {
	// 根据topic取到所有的分区
	consumer := base.KakfaConsumerClient()
	partitionList, err := consumer.Partitions(tailJob.topic)
	if err != nil {
		panic("获取topic所有的分区失败" + err.Error())
	}
	logrus.Debug("topic: ", tailJob.topic)

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(tailJob.topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return err
		}
		// defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					// 直接发给ES
					if err != nil {
						fmt.Printf("unmarshal failed. err:%v\n", err)
					} else {
						ld := EsLogData{Topic: tailJob.topic, Data: string(msg.Value)}
						// 优化一下: 直接放到一个chan中
						SendToESChan(&ld)
					}
				case <-tailJob.ctx.Done():
					logrus.Info("退出kakfa监听, 从新监听新的topic：", tailJob.topic)
					return
				default:

				}
			}
			//for msg := range pc.Messages() {
			//	fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			//	// 直接发给ES
			//	if err != nil {
			//		fmt.Printf("unmarshal failed. err:%v\n", err)
			//		continue
			//	}
			//	ld := EsLogData{Topic: tailJob.topic, Data: string(msg.Value)}
			//	// 优化一下: 直接放到一个chan中
			//	SendToESChan(&ld)
			//
			//}

		}(pc)
	}
	return
}
