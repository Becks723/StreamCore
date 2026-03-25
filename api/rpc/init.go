package rpc

import (
	"fmt"

	"StreamCore/config"
	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/chat/chatservice"
	"StreamCore/kitex_gen/group/groupservice"
	"StreamCore/kitex_gen/interaction/interactionservice"
	"StreamCore/kitex_gen/social/socialservice"
	"StreamCore/kitex_gen/user/userservice"
	"StreamCore/kitex_gen/video/videoservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	userClient   userservice.Client
	videoClient  videoservice.Client
	iaClient     interactionservice.Client
	socialClient socialservice.Client
	chatClient   chatservice.Client
	groupClient  groupservice.Client
)

func Init() {
	initUserRPC()
	initVideoRPC()
	initSocialRPC()
	initInteractionRPC()
	initChatRPC()
	initGroupRPC()
}

func InitRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) {
	config := config.Instance()
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("InitRPCClient: error etcd.NewEtcdResolver: %w", err)
	}
	c, err := newClientFunc(serviceName,
		client.WithResolver(r),
		client.WithMuxConnection(constants.MuxConnection),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err != nil {
		return nil, fmt.Errorf("InitRPCClient: error newClientFunc: %w", err)
	}
	return &c, nil
}

func initRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) {
	return InitRPCClient(serviceName, newClientFunc)
}
