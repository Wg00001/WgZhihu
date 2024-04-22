package interceptors

import (
	"WgZhihu/pkg/xcode"
	"context"

	"google.golang.org/grpc"
)

/*
ServerErrorInterceptor
拦截器，用于rpc的自定义错误码
由于rpc无法正常使用错误码，所以需要拦截器
*/
func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//获取正常业务处理对应的rpc方法处理的结果
		resp, err = handler(ctx, req)
		//将正常rpc返回的error（自己定义的业务错误码）转换成grpc的错误码
		return resp, xcode.FromError(err).Err()
	}
}
