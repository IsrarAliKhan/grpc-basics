package middleware

import (
	"context"
	"item/server/auth"
	"log"

	"google.golang.org/grpc"
)

func Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("---> unary interceptor: ", info.FullMethod)

		if info.FullMethod != "/proto.Auth/Login" {
			err := auth.Authorize(ctx, info.FullMethod)
			if err != nil {
				return nil, err
			}
		}

		return handler(ctx, req)
	}
}

func Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("---> stream interceptor: ", info.FullMethod)

		err := auth.Authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}
