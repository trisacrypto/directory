package client_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/shibukawa/configdir"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/client"
	"gopkg.in/yaml.v2"
)

const profileYAML = "profiles.yaml"

// Test that DefaultProfiles returns a Profiles object with some valid profiles.
func TestDefaultProfiles(t *testing.T) {
	profiles := client.DefaultProfiles()
	require.NotNil(t, profiles)
	require.True(t, strings.HasPrefix(profiles.Version, "v"))
	require.NotEmpty(t, profiles.Profiles)
	require.Contains(t, profiles.Profiles, profiles.Active)
}

// Test that GetProfilesFolder finds the correct folder.
func TestGetProfilesFolder(t *testing.T) {
	// Should default to the local path if it contains a profiles config
	f, err := os.Create(profileYAML)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	defer os.RemoveAll(profileYAML)
	folder, err := client.GetProfilesFolder()
	require.NoError(t, err)
	_, path, _, ok := runtime.Caller(0)
	require.True(t, ok)
	require.Equal(t, filepath.Dir(path), folder.Path)
}

// Test that Load correctly reads the profiles from disk.
func TestLoad(t *testing.T) {
	defer os.RemoveAll(profileYAML)

	// Config is empty
	f, err := os.Create(profileYAML)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	_, err = client.Load()
	require.Error(t, err)
	require.NoError(t, os.RemoveAll(profileYAML))

	// Config is badly formatted
	f, err = os.Create(profileYAML)
	require.NoError(t, err)
	_, err = f.WriteString("bad")
	require.NoError(t, err)
	require.NoError(t, f.Close())
	_, err = client.Load()
	require.Error(t, err)
	require.NoError(t, os.RemoveAll(profileYAML))

	// Config contains wrong version
	f, err = os.Create(profileYAML)
	require.NoError(t, err)
	wrongVersion := &client.Profiles{
		Version: "v2",
	}
	data, err := yaml.Marshal(wrongVersion)
	require.NoError(t, err)
	_, err = f.Write(data)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	_, err = client.Load()
	require.Error(t, err)
	require.NoError(t, os.RemoveAll(profileYAML))

	// Config is valid
	f, err = os.Create(profileYAML)
	require.NoError(t, err)
	expected := client.DefaultProfiles()
	data, err = yaml.Marshal(expected)
	require.NoError(t, err)
	_, err = f.Write(data)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	// Load should return the default profiles
	profiles, err := client.Load()
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
	defer os.RemoveAll(profileYAML)
	profiles := client.DefaultProfiles()
	f, err := os.Create(profileYAML)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	// Profile name does not exist
	err = profiles.SetActive("does-not-exist")
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

	// Active profile should be changed and new config written to disk
	err = profiles.SetActive(name)
	require.NoError(t, err)
	require.Equal(t, name, profiles.Active)
	data, err := os.ReadFile(profileYAML)
	require.NoError(t, err)
	actual := &client.Profiles{}
	err = yaml.Unmarshal(data, actual)
	require.NoError(t, err)
	require.Equal(t, profiles, actual)
}

// Test that Save writes the profiles to the specified folder.
func TestSave(t *testing.T) {
	defer os.RemoveAll(profileYAML)
	profiles := client.DefaultProfiles()
	f, err := os.Create(profileYAML)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	// Profiles should be written to the local config if not specified.
	err = profiles.Save(nil)
	require.NoError(t, err)
	actual, err := client.Load()
	require.NoError(t, err)
	require.Equal(t, profiles, actual)

	// Profiles should be written to the specified config.
	tmp, err := ioutil.TempDir("", "config-*")
	defer os.RemoveAll(tmp)
	require.NoError(t, err)
	folder := &configdir.Config{
		Path: tmp,
		Type: configdir.Local,
	}
	err = profiles.Save(folder)
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(tmp, profileYAML))
	require.NoError(t, err)
	actual = &client.Profiles{}
	err = yaml.Unmarshal(data, actual)
	require.NoError(t, err)
	require.Equal(t, profiles, actual)
}

// Test that Install creates default profiles and saves them to disk.
func TestInstall(t *testing.T) {
	defer os.RemoveAll(profileYAML)
	f, err := os.Create(profileYAML)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	err = client.Install()
	require.NoError(t, err)
	data, err := os.ReadFile(profileYAML)
	require.NoError(t, err)
	actual := &client.Profiles{}
	err = yaml.Unmarshal(data, actual)
	require.NoError(t, err)
	require.Equal(t, client.DefaultProfiles(), actual)
}
