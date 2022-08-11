package sectigo_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

func TestConfigValidation(t *testing.T) {
	conf := sectigo.Config{
		Profile: "invalid profile",
	}
	require.ErrorContains(t, conf.Validate(), fmt.Sprintf("%q is not a valid Sectigo profile name", conf.Profile))

	conf.Profile = sectigo.AllProfiles()[0]
	require.NoError(t, conf.Validate())
}
