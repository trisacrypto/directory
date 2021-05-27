package models

//go:generate protoc -I=../../../../../trisa/proto -I=../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/gds/models/v1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/gds/models/v1 gds/models/v1/models.proto
