/*
Package pkg describes the TRISA Global Directory Go reference package.
*/
package pkg

import "fmt"

// Version component constants for the current build.
const (
	VersionMajor = 0
	VersionMinor = 2
	VersionPatch = 1
)

// Version returns the semantic version for the current build.
func Version() string {
	if VersionPatch > 0 {
		return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	}
	return fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
}
