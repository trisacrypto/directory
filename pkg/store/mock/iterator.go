package mock

import (
	"github.com/trisacrypto/directory/pkg/store/iterator"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func BadDirectoryIterator(err error) iterator.DirectoryIterator {
	return &MockIterator{err: err}
}

// TODO: implement mock iteration for testing
type MockIterator struct {
	err error
}

func (m *MockIterator) Next() bool   { return false }
func (m *MockIterator) Prev() bool   { return false }
func (m *MockIterator) Error() error { return m.err }
func (m *MockIterator) Release()     {}

func (m *MockIterator) Id() string                { return "" }
func (m *MockIterator) VASP() (*pb.VASP, error)   { return nil, nil }
func (m *MockIterator) All() ([]*pb.VASP, error)  { return nil, nil }
func (m *MockIterator) SeekId(vaspID string) bool { return false }
