/*
Package pkg describes the TRISA Global Directory Go reference package.
*/
package pkg

import "fmt"

// Version component constants for the current build.
const (
	VersionMajor         = 1
	VersionMinor         = 1
	VersionPatch         = 0
	VersionReleaseLevel  = "alpha"
	VersionReleaseNumber = 2
)

// Version returns the semantic version for the current build.
func Version() string {
	var versionCore string
	if VersionPatch > 0 {
		versionCore = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	} else {
		versionCore = fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
	}

	if VersionReleaseLevel != "" {
		if VersionReleaseNumber > 0 {
			return fmt.Sprintf("%s-%s.%d", versionCore, VersionReleaseLevel, VersionReleaseNumber)
		}
		return fmt.Sprintf("%s-%s", versionCore, VersionReleaseLevel)
	}
	return versionCore
}
