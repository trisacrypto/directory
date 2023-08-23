package activity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/utils/activity"
)

func TestPublish(t *testing.T) {
	// Test using the global methods to publish activity events
	activity.Lookup().VASP(uuid.New()).Add()
}
