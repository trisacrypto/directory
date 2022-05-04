package members

//go:generate protoc -I=$GOPATH/src/github.com/trisacrypto/trisa/proto -I=$GOPATH/src/github.com/trisacrypto/directory/proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/gds/members/v1alpha1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/gds/members/v1alpha1 gds/members/v1alpha1/members.proto
