package members

//go:generate protoc -I=../../../../../trisa/proto -I=../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/gds/members/v1alpha1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/gds/members/v1alpha1 gds/members/v1alpha1/members.proto
