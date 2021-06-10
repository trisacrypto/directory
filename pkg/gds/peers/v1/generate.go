package peers

//go:generate protoc -I=../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/gds/peers/v1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/gds/peers/v1 gds/peers/v1/peers.proto
