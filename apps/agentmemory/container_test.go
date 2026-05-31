package main

import (
	"testing"

	"github.com/kalenarndt/containers/testhelpers"
)

func Test(t *testing.T) {
	image := testhelpers.GetTestImage("ghcr.io/kalenarndt/agentmemory:rolling")
	testhelpers.TestHTTPEndpoint(t, image, testhelpers.HTTPTestConfig{
		Port:       "3111",
		Path:       "/agentmemory/livez",
		StatusCode: 200,
	}, nil)
}
