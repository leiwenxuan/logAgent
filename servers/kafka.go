package servers

import (
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

func init() {
	logDataChan = make(chan *logData, 100000)
}

func SendToKafka() {
	for {
		select {
		case msgChan := <-logDataChan:
			msg := &sarama.ProducerMessage{
				Topic: msgChan.topic,
				Value: sarama.StringEncoder(msgChan.msg),
			}
			client := base.KakfaClient()
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Error("Send msg failed, err: ", err)
			}
			logrus.Debug("发送消息成功： ", pid, offset)

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
