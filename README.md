# HWMonit with gRPC and http API





### Update ProtoBuf

> API变动才需要执行一下步骤，更新go文件

```shell
cd network
protoc --go_out=plugins=grpc:. ./service.proto

# 编译google.api
protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/*.proto

#编译service.proto为service.pb.proto
protoc -I . --go_out=plugins=grpc,Mgoogle/api/annotations.proto=google/api:. ./service.proto

#编译service.proto为service.pb.gw.proto
protoc --grpc-gateway_out=logtostderr=true:. ./service.proto
```