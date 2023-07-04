package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionBetween(t *testing.T) {
	assert.True(t, VersionBetween("1", "2", "3"))
	assert.True(t, VersionBetween("1", "2", ""))
	assert.True(t, VersionBetween("", "2", "3"))
	assert.True(t, VersionBetween("", "2", ""))
	assert.True(t, VersionBetween("1", "", ""))
	assert.True(t, VersionBetween("", "", "3"))
	assert.True(t, VersionBetween("1", "", "3"))

	assert.True(t, VersionBetween("2", "2", "2"))
	assert.True(t, VersionBetween("2", "2", ""))
	assert.True(t, VersionBetween("", "2", "2"))
}

func testVersionSatifiesRange(t *testing.T, v, rng string) {
	satisfiesRange, err := VersionSatisfiesRange(v, rng)
	assert.Nil(t, err)
	assert.True(t, satisfiesRange)
}

func testNotVersionSatifiesRange(t *testing.T, v, rng string) {
	satisfiesRange, err := VersionSatisfiesRange(v, rng)
	assert.Nil(t, err)
	assert.False(t, satisfiesRange)
}

func testInvalidVersion(t *testing.T, v, rng string) {
	satisfiesRange, _ := VersionSatisfiesRange(v, rng)
	assert.False(t, satisfiesRange)
}

func TestVersionSatifiesRange(t *testing.T) {
	testVersionSatifiesRange(t, "v1.0.0", "=1.0.0")
	testVersionSatifiesRange(t, "1.0.0", "!2.0.0")
	testVersionSatifiesRange(t, "v1.0.2", ">1.0.1 <1.0.3")
	testVersionSatifiesRange(t, "1.0.0", "<1.0.1 || >1.0.3")
	testVersionSatifiesRange(t, "v1.0.4", "<1.0.1 || >1.0.3")
	testVersionSatifiesRange(t, "v1.0.0", "=v1.0.0")
	testVersionSatifiesRange(t, "1.0.0", "!v2.0.0")
	testVersionSatifiesRange(t, "v1.0.2", ">v1.0.1 <v1.0.3")
	testVersionSatifiesRange(t, "1.0.0", "<v1.0.1 || >v1.0.3")
	testVersionSatifiesRange(t, "v1.0.4", "<v1.0.1 || >v1.0.3")

	testVersionSatifiesRange(t, "v1.0.0-ranger11", "=1.0.0-ranger11")
	testVersionSatifiesRange(t, "1.0.0-ranger11", "!1.0.0-ranger12")
	testVersionSatifiesRange(t, "v1.0.0-ranger2", ">1.0.0-ranger1 <1.0.0-ranger3")
	testVersionSatifiesRange(t, "1.0.0-ranger1", "<1.0.0-ranger2 || >1.0.0-ranger4")
	testVersionSatifiesRange(t, "v1.0.0-ranger5", "<1.0.0-ranger2 || >1.0.0-ranger4")
	testVersionSatifiesRange(t, "v1.0.0-ranger11", "=v1.0.0-ranger11")
	testVersionSatifiesRange(t, "1.0.0-ranger11", "!v1.0.0-ranger12")
	testVersionSatifiesRange(t, "v1.0.0-ranger2", ">v1.0.0-ranger1 <v1.0.0-ranger3")
	testVersionSatifiesRange(t, "1.0.0-ranger1", "<v1.0.0-ranger2 || >v1.0.0-ranger4")
	testVersionSatifiesRange(t, "v1.0.0-ranger5", "<v1.0.0-ranger2 || >v1.0.0-ranger4")

	testVersionSatifiesRange(t, "v1.0.0-pre1-ranger11", "=1.0.0-pre1-ranger11")
	testVersionSatifiesRange(t, "1.0.0-pre1-ranger11", "!1.0.0-pre1-ranger12")
	testVersionSatifiesRange(t, "v1.0.0-pre1-ranger2", ">1.0.0-pre1-ranger1 <1.0.0-pre1-ranger3")
	testVersionSatifiesRange(t, "1.0.0-pre1-ranger1", "<1.0.0-pre1-ranger2 || >1.0.0-pre1-ranger4")
	testVersionSatifiesRange(t, "v1.0.0-pre1-ranger5", "<1.0.0-pre1-ranger2 || >1.0.0-pre1-ranger4")
	testVersionSatifiesRange(t, "v1.0.0-pre1-ranger11", "=v1.0.0-pre1-ranger11")
	testVersionSatifiesRange(t, "1.0.0-pre1-ranger11", "!v1.0.0-pre1-ranger12")
	testVersionSatifiesRange(t, "v1.0.0-pre1-ranger2", ">v1.0.0-pre1-ranger1 <v1.0.0-pre1-ranger3")
	testVersionSatifiesRange(t, "1.0.0-pre1-ranger1", "<v1.0.0-pre1-ranger2 || >v1.0.0-pre1-ranger4")
	testVersionSatifiesRange(t, "v1.0.0-pre1-ranger5", "<v1.0.0-pre1-ranger2 || >v1.0.0-pre1-ranger4")

	testVersionSatifiesRange(t, "v1.0.0-pre11-ranger1", "=1.0.0-pre11-ranger1")
	testVersionSatifiesRange(t, "1.0.0-pre11-ranger1", "!1.0.0-pre12-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre2-ranger1", ">1.0.0-pre1-ranger1 <1.0.0-pre3-ranger1")
	testVersionSatifiesRange(t, "1.0.0-pre1-ranger1", "<1.0.0-pre2-ranger1 || >1.0.0-pre4-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre5-ranger1", "<1.0.0-pre2-ranger1 || >1.0.0-pre4-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre11-ranger1", "=v1.0.0-pre11-ranger1")
	testVersionSatifiesRange(t, "1.0.0-pre11-ranger1", "!v1.0.0-pre12-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre2-ranger1", ">v1.0.0-pre1-ranger1 <v1.0.0-pre3-ranger1")
	testVersionSatifiesRange(t, "1.0.0-pre1-ranger1", "<v1.0.0-pre2-ranger1 || >v1.0.0-pre4-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre5-ranger1", "<v1.0.0-pre2-ranger1 || >v1.0.0-pre4-ranger1")

	testVersionSatifiesRange(t, "v1.0.0-pre11-ranger1", "=1.0.0-pre11-ranger1")
	testVersionSatifiesRange(t, "1.0.0-pre11-ranger1", "!1.0.0-pre12-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre2-ranger1", ">1.0.0-pre1-ranger1 <1.0.0-pre3-ranger1")
	testVersionSatifiesRange(t, "1.0.0-pre1-ranger1", "<1.0.0-pre2-ranger1 || >1.0.0-pre4-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre5-ranger1", "<1.0.0-pre2-ranger1 || >1.0.0-pre4-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre11-ranger1", "=v1.0.0-pre11-ranger1")
	testVersionSatifiesRange(t, "1.0.0-pre11-ranger1", "!v1.0.0-pre12-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre2-ranger1", ">v1.0.0-pre1-ranger1 <v1.0.0-pre3-ranger1")
	testVersionSatifiesRange(t, "1.0.0-pre1-ranger1", "<v1.0.0-pre2-ranger1 || >v1.0.0-pre4-ranger1")
	testVersionSatifiesRange(t, "v1.0.0-pre5-ranger1", "<v1.0.0-pre2-ranger1 || >v1.0.0-pre4-ranger1")

	testVersionSatifiesRange(t, "v1.0.0-pre2-ranger1", ">1.0.0-pre1-ranger2")
	testVersionSatifiesRange(t, "v1.0.0-pre2-ranger1", "<1.0.0")
	testVersionSatifiesRange(t, "v1.0.0-pre2-ranger1", ">v1.0.0-pre1-ranger2")
	testVersionSatifiesRange(t, "v1.0.0-pre2-ranger1", "<v1.0.0")

	testNotVersionSatifiesRange(t, "v1.0.0-ranger12", "=1.0.0-ranger11")
	testNotVersionSatifiesRange(t, "1.0.0-ranger12", "!1.0.0-ranger12")
	testNotVersionSatifiesRange(t, "v1.0.0-ranger5", ">1.0.0-ranger1 <1.0.0-ranger3")
	testNotVersionSatifiesRange(t, "1.0.0-ranger3", "<1.0.0-ranger2 || >1.0.0-ranger4")
	testNotVersionSatifiesRange(t, "v1.0.0-ranger12", "=v1.0.0-ranger11")
	testNotVersionSatifiesRange(t, "1.0.0-ranger12", "!v1.0.0-ranger12")
	testNotVersionSatifiesRange(t, "v1.0.0-ranger5", ">v1.0.0-ranger1 <v1.0.0-ranger3")
	testNotVersionSatifiesRange(t, "1.0.0-ranger3", "<v1.0.0-ranger2 || >v1.0.0-ranger4")

	testInvalidVersion(t, "versionInvalid-1.0", "versionInvalid-1.0")
	testInvalidVersion(t, "versionInvalid-1.0", "=versionInvalid-1.0")
	testInvalidVersion(t, "versionInvalid-1.0", "<versionInvalid-1.0")
	testInvalidVersion(t, "versionInvalid-1.0", "<=versionInvalid-1.0")
	testInvalidVersion(t, "versionInvalid-1.0", ">versionInvalid-1.0")
	testInvalidVersion(t, "versionInvalid-1.0", ">=versionInvalid-1.0")

	testInvalidVersion(t, "v1.0.0-validVersion", "versionInvalid-1.0")
	testInvalidVersion(t, "v1.0.0-validVersion", "=versionInvalid-1.0")
	testInvalidVersion(t, "v1.0.0-validVersion", ">versionInvalid-1.0")
	testInvalidVersion(t, "v1.0.0-validVersion", ">=versionInvalid-1.0")
	testInvalidVersion(t, "v1.0.0-validVersion", "<versionInvalid-1.0")
	testInvalidVersion(t, "v1.0.0-validVersion", "<=versionInvalid-1.0")

	testInvalidVersion(t, "versionInvalid-1.0", "v1.0.0-validVersion")
	testInvalidVersion(t, "versionInvalid-1.0", "=v1.0.0-validVersion")
	testInvalidVersion(t, "versionInvalid-1.0", ">v1.0.0-validVersion")
	testInvalidVersion(t, "versionInvalid-1.0", ">=v1.0.0-validVersion")
	testInvalidVersion(t, "versionInvalid-1.0", "<v1.0.0-validVersion")
	testInvalidVersion(t, "versionInvalid-1.0", "<=v1.0.0-validVersion")

}
