/*
Package pkg describes the TRISA Global Directory Go reference package.
*/
package pkg

import "fmt"

// Version component constants for the current build.
const (
	VersionMajor         = 1
	VersionMinor         = 8
	VersionPatch         = 0
	VersionReleaseLevel  = ""
	VersionReleaseNumber = 16
)

// Set the GitVersion via -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$(git rev-parse --short HEAD)'"
var GitVersion string

// Version returns the semantic version for the current build.
func Version() string {
	var versionCore string
	if VersionPatch > 0 || VersionReleaseLevel != "" {
		versionCore = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	} else {
		versionCore = fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
	}

	if VersionReleaseLevel != "" {
		if VersionReleaseNumber > 0 {
			versionCore = fmt.Sprintf("%s-%s.%d", versionCore, VersionReleaseLevel, VersionReleaseNumber)
		} else {
			versionCore = fmt.Sprintf("%s-%s", versionCore, VersionReleaseLevel)
		}
	}

	if GitVersion != "" {
		versionCore = fmt.Sprintf("%s (%s)", versionCore, GitVersion)
	}

	return versionCore
}
