package interceptors

import (
	"WgZhihu/pkg/xcode"
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ClientErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			//将grpc的status转成自定义的xcode
			xc := xcode.GrpcStatusToXCode(grpcStatus)
			err = errors.WithMessage(xc, grpcStatus.Message())
		}

		return err
	}
}
