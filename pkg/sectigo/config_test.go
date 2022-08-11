package sectigo_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

func TestConfigValidation(t *testing.T) {
	conf := sectigo.Config{
		Profile: "invalid profile",
	}
	require.EqualError(t, conf.Validate(), fmt.Sprintf("%q is not a valid Sectigo profile name, specify one of %s", conf.Profile, strings.Join(sectigo.AllProfiles(), ", ")))

	conf.Profile = sectigo.AllProfiles()[0]
	require.NoError(t, conf.Validate())
}
