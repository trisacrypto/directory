package client

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shibukawa/configdir"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// Profiles are stored in a system-specific configuration location e.g.
// ~/Library/Application Support/rotational/gds on OS X or ~/.config/rotational/gds on
// Linux. The profiles make it easy to switch between client configurations to connect
// to trisatest.net or vaspdirectory.net. The profiles have a user-supplied name for
// easy configuration and the profiles are populated with reasonable defaults.
//
// At most one profile can be marked as "active", this profile is treated as the default
// profile if a specific profile is not used. If multiple profiles are marked as active
// the "first" profile marked as active is used (with no guaranteed ordering).
type Profiles map[string]*Profile

var cfgd = configdir.New("rotational", "gds")

const (
	profileYAML = "profiles.yaml"
)

// ProfilesPath returns the location on disk where the profiles are stored. If no
// profiles are located then an error is returned.
func ProfilesPath() (string, error) {
	folder := cfgd.QueryFolderContainsFile(profileYAML)
	if folder != nil {
		return filepath.Join(folder.Path, profileYAML), nil
	}
	return "", errors.New("no profiles are available")
}

// Load the profiles from disk if they're available.
func Load() (p Profiles, err error) {
	folder := cfgd.QueryFolderContainsFile(profileYAML)
	if folder != nil {
		var data []byte
		if data, err = folder.ReadFile(profileYAML); err != nil {
			return nil, fmt.Errorf("could not read %s: %s", profileYAML, err)
		}

		p = make(Profiles)
		if err = yaml.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("could not unmarshal profiles: %s", err)
		}

		return p, nil
	}
	return nil, errors.New("no profiles are available")
}

// LoadActive is a shorthand for Load() then Active() and finally Update()
func LoadActive(c *cli.Context) (p *Profile, err error) {
	var profiles Profiles
	if profiles, err = Load(); err != nil {
		return nil, err
	}

	if p, err = profiles.Active(c.String("profile")); err != nil {
		return nil, err
	}

	if err = p.Update(c); err != nil {
		return nil, err
	}

	return p, nil
}

// Active returns the profile with the specified name or the active profile if no name
// is specified. If multiple profiles are marked active it returns the first active
// profile with no ordering gurantees. If no profiles are marked active and there is
// one profile, that profile is returned, otherwise an error is returned.
func (p Profiles) Active(name string) (_ *Profile, err error) {
	if name != "" {
		profile, ok := p[name]
		if !ok {
			return nil, fmt.Errorf("no profile found named %q", name)
		}
		return profile, nil
	}

	for _, profile := range p {
		if profile.Active {
			return profile, nil
		}
	}

	if len(p) == 1 {
		for _, profile := range p {
			return profile, nil
		}
	}

	return nil, fmt.Errorf("no active profile found in %d profiles", len(p))
}

// SetActive marks the profile with the specified name as active.
func (p Profiles) SetActive(name string) (err error) {
	if _, ok := p[name]; !ok {
		return fmt.Errorf("no profile named %q found", name)
	}

	// Mark all profiles inactive to ensure only one profile is active at a time
	for _, profile := range p {
		profile.Active = false
	}

	// Mark the specified profile as active
	p[name].Active = true

	// Save the profiles back to disk to ensure the activation takes effect
	if err = p.Save(); err != nil {
		return fmt.Errorf("could not save active profile: %s", err)
	}
	return nil
}

// Save the profiles to disk in the configuration directory.
func (p Profiles) Save() (err error) {
	folders := cfgd.QueryFolders(configdir.Global)
	if len(folders) == 0 {
		return errors.New("could not find user configuration directory")
	}

	var data []byte
	if data, err = yaml.Marshal(p); err != nil {
		return fmt.Errorf("could not marshal profiles: %s", err)
	}

	// Save the configuration to the first folder
	if err = folders[0].WriteFile(profileYAML, data); err != nil {
		return fmt.Errorf("could not write profiles to disk: %s", err)
	}

	return nil
}

// Install creates a default profile and saves it to disk, overwriting the previous contents.
func Install() (err error) {
	profiles := make(Profiles)
	profiles["production"] = &Profile{
		Directory: &DirectoryProfile{
			Endpoint: "api.vaspdirectory.net:443",
		},
		Admin: &AdminProfile{
			Endpoint: "https://api.admin.vaspdirectory.net",
		},
	}
	profiles["testnet"] = &Profile{
		Directory: &DirectoryProfile{
			Endpoint: "api.trisatest.net:443",
		},
		Admin: &AdminProfile{
			Endpoint: "https://api.admin.trisatest.net",
		},
		Active: true,
	}
	profiles["localhost"] = &Profile{
		Directory: &DirectoryProfile{
			Endpoint: "localhost:4433",
			Insecure: true,
		},
		Admin: &AdminProfile{
			Endpoint: "http://localhost:4434",
		},
		DatabaseURL: os.Getenv("GDS_DATABASE_URL"),
	}

	return profiles.Save()
}
