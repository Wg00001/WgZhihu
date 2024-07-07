package logic

import (
	"WgZhihu/application/like/mq/internal/svc"
	"context"
	"fmt"
	"github.com/zeromicro/go-queue/kq"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Consume(key, val string) error {
	fmt.Printf("get key: %s val: %s\n", key, val)
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	//将kafka注册成service，通过service start方式启动服务
	return []service.Service{
		//将实现了Consume方法的接口，传给kq.MustNewQueue()
		//返回一个service。在main中，通过serviceGroup来启动服务
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewThumbupLogic(ctx, svcCtx)),
	}
}
