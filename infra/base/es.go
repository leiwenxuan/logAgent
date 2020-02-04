package base

import (
	"github.com/olivere/elastic/v7"

	"gitee.com/lwx0416/logAgent/infra"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

//dbx 数据库实例
var clientES *elastic.Client

func ESClient() *elastic.Client {
	Check(clientES)
	return clientES
}

//ES starter，并且设置为全局
type ESStarter struct {
	infra.BaseStarter
}

func (s *ESStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()
	var (
		client *elastic.Client
		err    error
	)
	endpoints := conf.GetDefault("es.endpoints", "http://127.0.0.1:9200")
	log.Debug("es.endpoints", endpoints)
	//client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(endpoints))
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetHealthcheck(false), elastic.SetURL(endpoints))
	if err != nil {
		// Handle error
		panic(err)
	}

	clientES = client
	log.Debug("clientMongo etcd:", clientES)
	return
}
