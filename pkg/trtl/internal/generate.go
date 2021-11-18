package internal

//go:generate protoc -I=../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/directory/pkg/trtl/internal  trtl/internal/pagination.proto
