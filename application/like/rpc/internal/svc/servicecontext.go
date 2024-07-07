package svc

import (
	"WgZhihu/application/like/rpc/internal/config"
	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config config.Config
	//kafka发送者客户端
	KqPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
