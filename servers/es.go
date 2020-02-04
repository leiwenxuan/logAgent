package servers

import (
	"context"
	"fmt"
	"time"

	"gitee.com/lwx0416/logAgent/infra/base"

	"github.com/sirupsen/logrus"
)

type EsLogData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

var (
	esMsgChan chan *EsLogData
)

func InitEs(nums int) {

	conf := base.Props()
	chanSize := conf.GetIntDefault("es.chanSize", 1000000)

	fmt.Println("初始化ES, chanSize", chanSize)
	//chanSize := 100000
	esMsgChan = make(chan *EsLogData, chanSize)

	for i := 0; i < nums; i++ {
		go sendToES()
	}
}

func SendToESChan(msg *EsLogData) {
	fmt.Println(esMsgChan)
	esMsgChan <- msg
}

func sendToES() {
	for {
		select {
		case msg := <-esMsgChan:
			client := base.ESClient()
			put, err := client.Index().Index(msg.Topic).Type("xxx").BodyJson(msg).Do(context.Background())
			if err != nil {
				logrus.Error(err)
				continue
			}
			logrus.Printf("Indexed log %v to index %s, type %s\n", put.Id, put.Index, put.Type)
		default:
			time.Sleep(time.Second)
		}
	}
}
