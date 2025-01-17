package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testAscending(t *testing.T, versions []string) {
	for i, version := range versions {
		for j := i; j < len(versions); j++ {
			if j != i {
				assert.True(t, GreaterThan(versions[j], version))
				assert.False(t, GreaterThan(version, versions[j]))
			}
		}
	}
}

func TestGreaterThan(t *testing.T) {
	testAscending(t, []string{
		"v1.2.0",
		"v1.2.1",
		"v1.2.3",
		"v1.3.0",
		"v1.3.4",
		"v2.0.0",
	})

	testAscending(t, []string{
		"v0.1.0-ranger0",
		"v0.1.0-ranger1",
		"v0.1.0-ranger1.1",
		"v1.2.4-ranger6",
		"v1.2.4-ranger6.1",
		"v1.2.4-ranger7",
		"v1.2.4-ranger7.2",
		"v1.2.4-ranger7.3",
		"v1.2.4-ranger9.0",
		"v1.2.4-ranger10.10",
		"v1.2.4-ranger12.0",
		"v1.2.4-ranger12.5",
		"v1.2.4-ranger14",
		"v1.2.4-ranger15.10",
		"v1.3.0-ranger3",
		"v1.3.0-ranger4",
	})

	testAscending(t, []string{
		"0.0.1",
		"v0.45.0",
	})

	testAscending(t, []string{
		"0.0.1-a",
		"0.0.1-b",
		"0.0.1-c",
	})

	testAscending(t, []string{
		"0.0.1-pre1-alpha2",
		"0.0.1-pre1-alpha3",
		"0.0.1-pre1-beta1",
		"0.0.1-pre1-beta2.2",
		"0.0.1-pre1-beta11",
		"0.0.1-pre1-rc1",
		"0.0.1-pre1-rc1-1",
		"0.0.1-pre1",
		"0.0.1",
	})

	assert.False(t, GreaterThan("v1.0.0+test", "v1.0.0"))
	assert.False(t, GreaterThan("v1.0.0", "v1.0.0+test"))
}
