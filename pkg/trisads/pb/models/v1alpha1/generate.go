package models

//go:generate protoc -I=../../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/trisads/pb/models/v1alpha1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/directory/pkg/trisads/pb/models/v1alpha1 trisads/models/v1alpha1/models.proto trisads/models/v1alpha1/ca.proto
