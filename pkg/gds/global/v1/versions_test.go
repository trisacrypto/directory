package global_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/config"
	. "github.com/trisacrypto/directory/pkg/gds/global/v1"
)

func TestVersionManager(t *testing.T) {
	conf := config.ReplicaConfig{}

	// Check required settings
	_, err := New(conf)
	require.Error(t, err)

	conf.PID = 8
	_, err = New(conf)
	require.Error(t, err)

	conf.Region = "us-east-2c"
	vers1, err := New(conf)
	require.NoError(t, err)
	require.Equal(t, "8:us-east-2c", vers1.Owner)

	conf.BindAddr = "us2.vaspdirectory.net:443"
	vers1, err = New(conf)
	require.NoError(t, err)
	require.Equal(t, "8:us2.vaspdirectory.net", vers1.Owner)

	conf.Name = "mitchell"
	vers1, err = New(conf)
	require.NoError(t, err)
	require.Equal(t, "8:mitchell", vers1.Owner)

	// Check update system
	require.Error(t, vers1.Update(nil))

	// Expected object definition:
	// &Object{
	// 	Key:       "foo",
	// 	Namespace: "awesome",
	// 	Region:    "us-east-2c",
	// 	Owner:     "8:mitchell",
	// 	Version: &Version{
	// 		PID:     8,
	// 		Version: 1,
	// 		Region:  "us-east-2c",
	// 		Parent:  nil,
	// 	},
	// }

	obj := &Object{Key: "foo", Namespace: "awesome"}
	require.NoError(t, vers1.Update(obj))
	require.Equal(t, vers1.Region, obj.Region)
	require.Equal(t, vers1.Owner, obj.Owner)
	require.Equal(t, vers1.PID, obj.Version.Pid)
	require.Equal(t, uint64(1), obj.Version.Version)
	require.Equal(t, vers1.Region, obj.Version.Region)
	require.Empty(t, obj.Version.Parent)

	// Create a new remote versioner
	conf.PID = 13
	conf.Region = "europe-west-3"
	conf.Name = "jaques"
	vers2, err := New(conf)
	require.NoError(t, err)

	// Update the previous version
	require.NoError(t, vers2.Update(obj))

	// Expected object definition:
	// &Object{
	// 	Key:       "foo",
	// 	Namespace: "awesome",
	// 	Region:    "us-east-2c",
	// 	Owner:     "8:mitchell",
	// 	Version: &Version{
	// 		Pid:     13,
	// 		Version: 2,
	// 		Region:  "europe-west-3",
	// 		Parent: &Version{
	// 			Pid:     8,
	// 			Version: 1,
	// 			Region:  "us-east-2c",
	// 		},
	// 	},
	// }

	require.Equal(t, vers1.Region, obj.Region)
	require.Equal(t, vers1.Owner, obj.Owner)
	require.Equal(t, vers2.PID, obj.Version.Pid)
	require.Equal(t, uint64(2), obj.Version.Version)
	require.Equal(t, vers2.Region, obj.Version.Region)
	require.NotEmpty(t, obj.Version.Parent)
	require.False(t, obj.Version.Parent.IsZero())
	require.Equal(t, vers1.PID, obj.Version.Parent.Pid)
	require.Equal(t, uint64(1), obj.Version.Parent.Version)
	require.Equal(t, vers1.Region, obj.Version.Parent.Region)
}
