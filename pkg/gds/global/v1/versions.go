package global

import (
	"errors"
	"fmt"
	"net"

	"github.com/trisacrypto/directory/pkg/gds/config"
)

// New creates a new global.VersionManager for handling lamport scalar versions.
func New(conf config.ReplicaConfig) (v *VersionManager, err error) {
	// Make sure we don't create a VersionManager that is unable to do its job.
	if conf.PID == 0 || conf.Region == "" {
		return nil, errors.New("improperly configured: version manager requires PID and Region")
	}

	v = &VersionManager{PID: conf.PID, Region: conf.Region}

	// Compute the owner name
	if conf.Name != "" {
		// The common name
		v.Owner = fmt.Sprintf("%d:%s", conf.PID, conf.Name)
	} else {
		// Check to see if there is a domain name in the bindaddr
		var host string
		if host, _, err = net.SplitHostPort(conf.BindAddr); err == nil {
			v.Owner = fmt.Sprintf("%d:%s", conf.PID, host)
		} else {
			// The owner name is just the pid:region in the last case
			v.Owner = fmt.Sprintf("%d:%s", conf.PID, conf.Region)
		}
	}

	return v, nil
}

type VersionManager struct {
	PID    uint64
	Owner  string
	Region string
}

// Update the version of an object in place.
func (v VersionManager) Update(meta *Object) error {
	if meta == nil {
		return errors.New("cannot update version on empty object")
	}

	// Update the parent to the current version of the object.
	if meta.Version != nil && !meta.Version.IsZero() {
		meta.Version.Parent = &Version{
			Pid:     meta.Version.Pid,
			Version: meta.Version.Version,
			Region:  meta.Version.Region,
		}
	} else {
		// This is the first version of the object. Also set provenance on the object.
		meta.Version = &Version{}
		meta.Region = v.Region
		meta.Owner = v.Owner
	}

	// Update the version to the new version of the local version manager
	meta.Version.Pid = v.PID
	meta.Version.Version++
	meta.Version.Region = v.Region
	return nil
}

// IsZero determines if the version is zero valued (e.g. the PID and Version are zero).
// Note that zero-valuation does not check parent or region.
func (v *Version) IsZero() bool {
	return v.Pid == 0 && v.Version == 0
}
