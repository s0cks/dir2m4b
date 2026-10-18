//go:build release

package internal

import "fmt"

func GetVersion() string {
	if Version == "latest" {
		return fmt.Sprintf("%s-release", Version)
	}

	return Version
}
