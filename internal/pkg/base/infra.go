package base

import (
	"log"
	"sync"

	"StreamCore/api/rpc"
	"StreamCore/internal/pkg/base/infra"
	"StreamCore/internal/pkg/cache"
	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db"
	"StreamCore/internal/pkg/mq"
	"StreamCore/kitex_gen/group/groupservice"
	"github.com/elastic/go-elasticsearch/v8"
)

type InfraSet struct {
	Cache       *cache.CacheSet
	DB          *db.DatabaseSet
	ES          *elasticsearch.TypedClient
	MQ          *mq.MQSet
	GroupClient groupservice.Client
}

var (
	instance *InfraSet
	once     sync.Once
)

type Option func(*InfraSet)

func GetInfraSet(opt ...Option) *InfraSet {
	once.Do(func() {
		instance = &InfraSet{}
		for _, op := range opt {
			op(instance)
		}
	})
	return instance
}

func WithDB() Option {
	return func(s *InfraSet) {
		gdb, err := infra.InitMySQL()
		if err != nil {
			log.Fatal(err)
		}
		s.DB = db.NewDatabaseSet(gdb)
	}
}

func WithCache() Option {
	return func(s *InfraSet) {
		rdb, err := infra.InitRedis()
		if err != nil {
			log.Fatal(err)
		}
		s.Cache = cache.NewCacheSet(rdb)
	}
}

func WithES() Option {
	return func(s *InfraSet) {
		es, err := infra.InitElastic()
		if err != nil {
			log.Fatal(err)
		}
		s.ES = es
	}
}

func WithMQ() Option {
	return func(s *InfraSet) {
		conn, err := infra.InitRabbitMQ()
		if err != nil {
			log.Fatal(err)
		}
		s.MQ = mq.NewMQSet(conn)
	}
}

func WithGroupClient() Option {
	return func(s *InfraSet) {
		c, err := rpc.InitRPCClient(constants.GroupServiceName, groupservice.NewClient)
		if err != nil {
			log.Fatal(err)
		}
		s.GroupClient = *c
	}
}
