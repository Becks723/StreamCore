package main

import (
	"log"
	"net"

	"StreamCore/config"
	groupimpl "StreamCore/internal/group"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/group/groupservice"
	"StreamCore/pkg/util"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	infra       *base.InfraSet
	serviceName = constants.GroupServiceName
	logPrefix   = "[group]"
)

func init() {
	config.Init(serviceName)
	infra = base.GetInfraSet(
		base.WithDB(),
		base.WithCache(),
		base.WithES(),
		base.WithMQ())
}

func main() {
	config := config.Instance()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		log.Fatalf("%s NewEtcdRegistry error: %v", logPrefix, err)
	}
	listenAddr, ok := util.GetAvailablePort(config.Service.AddrList)
	if !ok {
		log.Fatalf("%s no port available", logPrefix)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatalf("%s ResolveTCPAddr error: %v", logPrefix, err)
	}

	svr := groupservice.NewServer(
		groupimpl.NewGroupHandler(infra),
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
