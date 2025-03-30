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

func init() {
	// Initialize default config dir for profiles
	cfgd = configdir.New("rotational", "gds")
	cfgd.LocalPath, _ = filepath.Abs(".")
}

var cfgd configdir.ConfigDir

// Profiles are stored in one of three locations, and are searched for in the following
// order:
// 1. Current directory (".")
// 2. User level directory (~/Library/Application Support/rotational/gds on OS X
// or ~/.config/rotational/gds on Linux)
// 3. System level directory (/Library/Application Suport/rotational/gds on OS X
// or /etc/xdg on Linux)
// If no profile config is found, one will be created in the first available directory
// based on the search order: Current directory -> User -> System. This allows the user
// to easily override the user or system config by creating a profiles.yaml in the CWD.
//
// The profiles make it easy to switch between client configurations to connect
// to testnet.directory or trisa.directory. The profiles have a user-supplied name for
// easy configuration and the profiles are populated with reasonable defaults.
//
// At most one profile is considered "active", this profile is treated as the default
// profile if a specific profile is not used.
type Profiles struct {
	Version  string              `yaml:"version"`
	Active   string              `yaml:"active"`
	Profiles map[string]*Profile `yaml:"profiles"`
}

const (
	ProfileYAML    = "profiles.yaml"
	ProfileVersion = "v1"
)

func DefaultProfiles() *Profiles {
	return &Profiles{
		Version: ProfileVersion,
		Active:  "localhost",
		Profiles: map[string]*Profile{
			"production": {
				Directory: &DirectoryProfile{
					Endpoint: "api.trisa.directory:443",
					Insecure: true,
				},
				Admin: &AdminProfile{
					Endpoint: "https://api.admin.trisa.directory",
				},
				Members: &MembersProfile{
					Endpoint: "members.trisa.directory:443",
					Insecure: true,
				},
				TrtlProfiles: []*TrtlProfile{
					{
						Endpoint: "trtl.us.trisa.directory:443",
						Insecure: true,
					},
				},
			},
			"testnet": {
				Directory: &DirectoryProfile{
					Endpoint: "api.testnet.directory:443",
					Insecure: true,
				},
				Admin: &AdminProfile{
					Endpoint: "https://api.admin.testnet.directory",
				},
				Members: &MembersProfile{
					Endpoint: "members.testnet.directory:443",
					Insecure: true,
				},
				TrtlProfiles: []*TrtlProfile{
					{
						Endpoint: "trtl.us.testnet.directory:443",
						Insecure: true,
					},
				},
			},
			"localhost": {
				Directory: &DirectoryProfile{
					Endpoint: "localhost:4433",
					Insecure: true,
				},
				Admin: &AdminProfile{
					Endpoint: "http://localhost:4434",
				},
				DatabaseURL: os.Getenv("GDS_DATABASE_URL"),
				Members: &MembersProfile{
					Endpoint: "localhost:4435",
					Insecure: true,
				},
				TrtlProfiles: []*TrtlProfile{
					{
						Endpoint: "localhost:4436",
						Insecure: true,
					},
				},
			},
		},
	}
}

// GetProfilesFolder returns a pointer to the folder where the profiles are stored. If
// no such folder is configured, it creates an empty config file in a suitable folder.
func GetProfilesFolder() (folder *configdir.Config, err error) {
	folder = cfgd.QueryFolderContainsFile(ProfileYAML)
	if folder == nil {
		// Search for an available folder to create the config file
		var folders []*configdir.Config
		folders = cfgd.QueryFolders(configdir.Global)
		if len(folders) == 0 {
			folders = cfgd.QueryFolders(configdir.System)
			if len(folders) == 0 {
				folders = cfgd.QueryFolders(configdir.Local)
				if len(folders) == 0 {
					return nil, errors.New("no suitable directory for config file")
				}
			}
		}

		// Create a new default config file under the directory we just located
		folder = folders[0]
		p := DefaultProfiles()
		if err = p.Save(folder); err != nil {
			return nil, fmt.Errorf("could not create new config file at %s: %s", filepath.Join(folder.Path, ProfileYAML), err)
		}
	}
	return folder, nil
}

// ProfilesPath returns the location on disk where the profiles are stored. If no
// profiles are located then an error is returned.
func ProfilesPath() (string, error) {
	folder, err := GetProfilesFolder()
	if err != nil {
		return "", fmt.Errorf("no profiles are available: %s", err)
	}
	return filepath.Join(folder.Path, ProfileYAML), nil
}

// Load the profiles from disk if they're available.
func Load() (p *Profiles, err error) {
	var folder *configdir.Config
	if folder, err = GetProfilesFolder(); err == nil {
		var data []byte
		if data, err = folder.ReadFile(ProfileYAML); err != nil {
			return nil, fmt.Errorf("could not read %s: %s", filepath.Join(folder.Path, ProfileYAML), err)
		}

		if err = yaml.Unmarshal(data, &p); err != nil {
			return nil, fmt.Errorf("could not unmarshal profiles: %s", err)
		}

		if p == nil {
			return nil, fmt.Errorf("profile config is empty: %s", filepath.Join(folder.Path, ProfileYAML))
		}

		if p.Version != ProfileVersion {
			return nil, fmt.Errorf("invalid profile version %s, expected %s", p.Version, ProfileVersion)
		}

		return p, nil
	}
	return nil, fmt.Errorf("no profiles are available: %s", err)
}

// LoadActive is a shorthand for Load() then GetActive() and finally Update()
func LoadActive(c *cli.Context) (p *Profile, err error) {
	var profiles *Profiles
	if profiles, err = Load(); err != nil {
		return nil, err
	}

	if p, err = profiles.GetActive(c.String("profile")); err != nil {
		return nil, err
	}

	if err = p.Update(c); err != nil {
		return nil, err
	}

	return p, nil
}

// GetActive returns the profile with the specified name or the active profile if no name
// is specified.
func (p *Profiles) GetActive(name string) (_ *Profile, err error) {
	if name != "" {
		profile, ok := p.Profiles[name]
		if !ok {
			return nil, fmt.Errorf("no profile found named %q", name)
		}
		return profile, nil
	}

	return p.Profiles[p.Active], nil
}

// SetActive marks the profile with the specified name as active.
func (p *Profiles) SetActive(name string) (err error) {
	if _, ok := p.Profiles[name]; !ok {
		return fmt.Errorf("no profile named %q found", name)
	}

	p.Active = name

	// Save the profiles back to disk to ensure the activation takes effect
	if err = p.Save(nil); err != nil {
		return fmt.Errorf("could not save active profile: %s", err)
	}
	return nil
}

// Save the profiles to disk in the specified configuration folder. If the configuration
// folder is nil, the configuration folder is located and created if it doesn't exist.
func (p *Profiles) Save(folder *configdir.Config) (err error) {
	if folder == nil {
		if folder, err = GetProfilesFolder(); err != nil {
			return fmt.Errorf("could not find profiles folder: %s", err)
		}
	}

	var data []byte
	if data, err = yaml.Marshal(p); err != nil {
		return fmt.Errorf("could not marshal profiles: %s", err)
	}

	// Save the configuration to the folder
	if err = folder.WriteFile(ProfileYAML, data); err != nil {
		return fmt.Errorf("could not write profiles to disk: %s", err)
	}

	return nil
}

// Install creates default profiles and saves them to disk, overwriting the previous contents.
func Install() (err error) {
	return DefaultProfiles().Save(nil)
}

// SetConfigDir is a helper utility to modify where the profiles package looks for the
// profiles.yaml file. This is generally used in tests but can also be used in
// environments where the default search path doesn't make sense.
func SetConfigDir(cd configdir.ConfigDir) {
	cfgd = cd
}
