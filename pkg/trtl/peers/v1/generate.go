package peers

//go:generate protoc -I=$GOPATH/src/github.com/trisacrypto/directory/proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/trtl/peers/v1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/trtl/peers/v1 trtl/peers/v1/peers.proto
