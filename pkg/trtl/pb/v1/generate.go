package pb

//go:generate protoc -I=../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/trtl/pb/v1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/trtl/pb/v1 trtl/v1/trtl.proto
