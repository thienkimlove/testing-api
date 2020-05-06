### Commands

```text
73  protoc -I shared/rpc/proto/listing/v1 -I shared/rpc/.third_party/googleapis --go_opt=paths=source_relative --go_out=plugins=grpc:custom_rpc/listing  shared/rpc/proto/listing/v1/service.proto
74  protoc -I shared/rpc/proto/listing/v1 -I shared/rpc/.third_party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:custom_rpc/listing  shared/rpc/proto/listing/v1/service.proto
75  git add . && git commit -m update && git push origin master
76  make generate
77  make run

```