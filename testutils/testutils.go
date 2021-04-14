package testutils

import (
	"testing"
)

func SkipCI(t *testing.T) {
	//if config.G_Properties.Profile == "ci" {
	t.Skip("Skipping testing in CI environment")
	//}
}
