package main

import (
	"flag"
	"fmt"

	"WgZhihu/application/like/rpc/internal/config"
	"WgZhihu/application/like/rpc/internal/server"
	"WgZhihu/application/like/rpc/internal/svc"
	svc2 "WgZhihu/application/like/rpc/types/service"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/like.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		svc2.RegisterLikeServer(grpcServer, server.NewLikeServer(ctx))

		//dev || test 模式下会注册反射服务。在配置文件中设置mode
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
