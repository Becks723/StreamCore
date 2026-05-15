package main

import (
	"log"
	"net"

	"StreamCore/config"
	aiimpl "StreamCore/internal/ai"
	"StreamCore/internal/pkg/ai/provider"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/ai/aiservice"
	"StreamCore/pkg/util"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	infra       *base.InfraSet
	serviceName = constants.AIServiceName
	logPrefix   = "[ai]"
)

func init() {
	config.Init(serviceName)
	provider.Init()
	infra = base.GetInfraSet(
		base.WithDB(),
		base.WithCache(),
		base.WithChatClient(),
	)
}

func main() {
	cfg := config.Instance()
	r, err := etcd.NewEtcdRegistry([]string{cfg.Etcd.Addr})
	if err != nil {
		log.Fatalf("%s NewEtcdRegistry error: %v", logPrefix, err)
	}
	listenAddr, ok := util.GetAvailablePort(cfg.Service.AddrList)
	if !ok {
		log.Fatalf("%s no port available", logPrefix)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatalf("%s ResolveTCPAddr error: %v", logPrefix, err)
	}

	svr := aiservice.NewServer(
		aiimpl.NewAIHandler(infra),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: serviceName,
		}),
		server.WithMuxTransport(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{
			MaxConnections: constants.MaxConnections,
			MaxQPS:         constants.MaxQPS,
		}),
	)
	if err = svr.Run(); err != nil {
		log.Fatalf("%s server.Run error: %v", logPrefix, err)
	}
}
