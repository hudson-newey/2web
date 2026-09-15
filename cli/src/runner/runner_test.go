package runner

import (
	"slices"
	"testing"

	"github.com/hudson-newey/2web-cli/src/packages"
)

// Every supported package manager must resolve to a runtime command; the
// switch in runtimeCommand once had empty case bodies for npm and pnpm, which
// made "2web serve" on ssr solutions fail with "Could not determine runtime".
func TestRuntimeCommandSupportsAllPackageManagers(t *testing.T) {
	tests := []struct {
		manager  packages.PackageManager
		expected []string
	}{
		{packages.Npm, []string{"node", "--experimental-strip-types"}},
		{packages.Pnpm, []string{"node", "--experimental-strip-types"}},
		{packages.Yarn, []string{"node", "--experimental-strip-types"}},
		{packages.Bun, []string{"bun"}},
		{packages.Deno, []string{"deno", "-A"}},
	}

	for _, test := range tests {
		command := runtimeCommand(test.manager)

		if !slices.Equal(command, test.expected) {
			t.Errorf(
				"runtimeCommand(%d) = %v, expected %v",
				test.manager,
				command,
				test.expected,
			)
		}
	}
}
