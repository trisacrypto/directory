package models

//go:generate protoc -I=$GOPATH/trisa/proto -I=$GOPATH/directory/proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/models/v1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/models/v1 gds/models/v1/models.proto
