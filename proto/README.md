# Installation

## Install go compiler

https://grpc.io/docs/languages/go/quickstart/

```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

# Compilation

## Compile to the auth api

```
protoc -Iproto --go_out=api-auth --go-grpc_out=api-auth --grpc-gateway_out=api-auth proto/*.proto
protoc -Iproto --go_out=api-user --go-grpc_out=api-user --grpc-gateway_out=api-user  proto/user.proto
protoc -Iproto --go_out=api-scrapper --go-grpc_out=api-scrapper --grpc-gateway_out=api-scrapper proto/*.proto
protoc -Iproto --go_out=api-media --go-grpc_out=api-media --grpc-gateway_out=api-media proto/*.proto
```

# GRPC/api_rest compatibility

https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/
https://grpc.io/blog/coreos/
https://googleapis.github.io/common-protos-java/1.14.0/apidocs/com/google/api/HttpRule.html
