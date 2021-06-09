package global

//go:generate protoc -I=../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/gds/global/v1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/gds/global/v1 gds/global/v1/global.proto
