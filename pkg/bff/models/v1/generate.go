package models

//go:generate protoc -I=$GOPATH/src/github.com/trisacrypto/trisa/proto -I=$GOPATH/src/github.com/trisacrypto/directory/proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/bff/models/v1 bff/models/v1/models.proto
