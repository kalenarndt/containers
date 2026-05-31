package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/postgres-init:rolling")
	testhelpers.TestFileExists(t, image, "/usr/local/bin/createdb", nil)
	testhelpers.TestFileExists(t, image, "/usr/local/bin/createuser", nil)
	testhelpers.TestFileExists(t, image, "/usr/local/bin/psql", nil)
	testhelpers.TestFileExists(t, image, "/usr/local/bin/pg_isready", nil)
}
