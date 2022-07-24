package internal

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

func GrpcMiddleware(ctx context.Context, r *http.Request) metadata.MD {
	md := make(map[string]string)

	// Extracts and saves in metatada the method of the request
	if method, ok := runtime.RPCMethod(ctx); ok {
		md["method"] = method
	}
	// Extracts and saves in metatada the method of the request
	if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
		md["pattern"] = pattern
	}

	return metadata.New(md)
}
