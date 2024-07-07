package main

import (
	"WgZhihu/application/like/mq/internal/config"
	"WgZhihu/application/like/mq/internal/logic"
	"WgZhihu/application/like/mq/internal/svc"
	"context"
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	//mq服务的启动需要service工具
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/like.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	//创建ServiceGroup后，将consumers注册到其中
	for _, mq := range logic.Consumers(ctx, svcCtx) {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}
