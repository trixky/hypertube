# Installation

## Install go compiler

https://grpc.io/docs/languages/go/quickstart/

# Compilation

## Compile to the auth api

```
protoc -Iproto --go_out=api-auth --go-grpc_out=api-auth  proto/*.proto
```

# GRPC/api_rest compatibility

https://grpc.io/blog/coreos/

https://googleapis.github.io/common-protos-java/1.14.0/apidocs/com/google/api/HttpRule.html