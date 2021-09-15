package admin

//go:generate protoc -I=../../../../../trisa/proto -I=../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/gds/admin/v1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/gds/admin/v1 gds/admin/v1/admin.proto
