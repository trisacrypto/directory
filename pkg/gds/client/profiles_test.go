package client_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/configdir"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/client"
	"gopkg.in/yaml.v2"
)

// Test that DefaultProfiles returns a Profiles object with some valid profiles.
func TestDefaultProfiles(t *testing.T) {
	profiles := client.DefaultProfiles()
	require.NotNil(t, profiles)
	require.Equal(t, profiles.Version, client.ProfileVersion)
	require.NotEmpty(t, profiles.Profiles)
	require.Contains(t, profiles.Profiles, profiles.Active)
}

// Test that GetProfilesFolder finds the correct folder.
func TestGetProfilesFolder(t *testing.T) {
	// Should default to the local path if it contains a profiles config
	path := makeTestConfigInDir(t, nil)

	folder, err := client.GetProfilesFolder()
	require.NoError(t, err)
	require.Equal(t, filepath.Dir(path), folder.Path)
}

// Test that Load correctly reads the profiles from disk.
func TestLoad(t *testing.T) {
	var err error

	// Config is empty
	makeTestConfigInDir(t, nil)
	_, err = client.Load()
	require.Error(t, err)

	// Config is badly formatted
	makeTestConfigInDir(t, "bad")
	_, err = client.Load()
	require.Error(t, err)

	// Config contains wrong version
	wrongVersion := &client.Profiles{
		Version: "v2",
	}
	makeTestConfigInDir(t, wrongVersion)
	_, err = client.Load()
	require.Error(t, err)

	// Config is valid - load should return the default profiles
	expected := client.DefaultProfiles()
	makeTestConfigInDir(t, expected)
	profiles, err := client.Load()
	require.NoError(t, err)
	require.Equal(t, expected, profiles)
}

// Test that GetActive returns the corect profile.
func TestGetActive(t *testing.T) {
	profiles := client.DefaultProfiles()

	// Default to active profile
	actual, err := profiles.GetActive("")
	require.NoError(t, err)
	require.Equal(t, profiles.Profiles[profiles.Active], actual)

	// Profile name does not exist
	_, err = profiles.GetActive("does-not-exist")
	require.Error(t, err)

	// Profile name specified
	actual, err = profiles.GetActive(profiles.Active)
	require.NoError(t, err)
	require.Equal(t, profiles.Profiles[profiles.Active], actual)
}

// Test that SetActive marks the specified profile as active and writes to disk.
func TestSetActive(t *testing.T) {
	profiles := client.DefaultProfiles()
	path := makeTestConfigInDir(t, profiles)

	// Profile name does not exist
	err := profiles.SetActive("does-not-exist")
	require.Error(t, err)

	// Find a profile that is not active
	name := ""
	for n := range profiles.Profiles {
		if n != profiles.Active {
			name = n
			break
		}
	}
	require.NotEmpty(t, name)

	// Active profile should be updated
	err = profiles.SetActive(name)
	require.NoError(t, err)
	require.Equal(t, name, profiles.Active)

	// Check that the change was written to the original profile
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	actual := &client.Profiles{}
	require.NoError(t, yaml.Unmarshal(data, actual), "could not unmarshal profiles")
	require.Equal(t, profiles, actual)
}

// Test that Save writes the profiles to the specified folder.
func TestSave(t *testing.T) {
	profiles := client.DefaultProfiles()
	makeTestConfigInDir(t, profiles)

	// Profiles should be written to the local config if not specified.
	err := profiles.Save(nil)
	require.NoError(t, err)
	actual, err := client.Load()
	require.NoError(t, err)
	require.Equal(t, profiles, actual)

	// Profiles should be written to the specified config.
	tmp, err := os.MkdirTemp("", "config-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	folder := &configdir.Config{
		Path: tmp,
		Type: configdir.Local,
	}
	err = profiles.Save(folder)
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(tmp, client.ProfileYAML))
	require.NoError(t, err)
	actual = &client.Profiles{}
	err = yaml.Unmarshal(data, actual)
	require.NoError(t, err)
	require.Equal(t, profiles, actual)
}

// Test that Install creates default profiles and saves them to disk.
func TestInstall(t *testing.T) {
	path := makeTestConfigInDir(t, nil)
	require.NoError(t, client.Install(), "could not install profiles, writing to a new file")

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	actual := &client.Profiles{}
	err = yaml.Unmarshal(data, actual)
	require.NoError(t, err)
	require.Equal(t, client.DefaultProfiles(), actual)
}

func setTestConfigDir(t *testing.T) configdir.ConfigDir {
	var err error
	cfgd := configdir.New("rotational", "gds-test")
	cfgd.LocalPath, err = os.MkdirTemp("", "config-*")
	require.NoError(t, err, "could not get tmp dir for configuration")
	client.SetConfigDir(cfgd)

	t.Cleanup(func() {
		os.RemoveAll(cfgd.LocalPath)
	})
	return cfgd
}

func makeTestConfigInDir(t *testing.T, fixture interface{}) string {
	cfgd := setTestConfigDir(t)
	path := filepath.Join(cfgd.LocalPath, client.ProfileYAML)

	if fixture != nil {
		data, err := yaml.Marshal(fixture)
		require.NoError(t, err, "could not marshal fixture")
		require.NoError(t, os.WriteFile(path, data, 0644), "could not write fixture to disk")
	} else {
		require.NoError(t, os.WriteFile(path, nil, 0644), "could not write empty fixture to disk")
	}

	return path
}
