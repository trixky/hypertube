package internal

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

func basic_middleware(ctx context.Context, _ *http.Request) metadata.MD {
	md := make(map[string]string)
	if method, ok := runtime.RPCMethod(ctx); ok {
		md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
	}
	if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
		md["pattern"] = pattern // /v1/example/login
	}
	return metadata.New(md)
}
