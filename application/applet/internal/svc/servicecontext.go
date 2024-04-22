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
	//客户端拦截器，用于将grpc传过来的status转换成我们自定义的错误码
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
