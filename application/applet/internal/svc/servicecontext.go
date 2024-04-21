package svc

import (
	"WgZhihu/application/applet/internal/config"
	"WgZhihu/application/user/rpc/user"
	"WgZhihu/pkg/interceptors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRPC  user.User
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	//conf := redis.RedisConf{
	//	Host: c.BizRedis.Host,
	//}
	//bizRedis, _ := redis.NewRedis(conf, redis.WithPass(c.BizRedis.Pass))
	return &ServiceContext{
		Config:   c,
		UserRPC:  user.NewUser(userRPC),
		BizRedis: redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
	}
}
