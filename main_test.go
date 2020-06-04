package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTag(t *testing.T) {
	version := "0.1.0"
	major, minor, patch, buildNumber := parseTag(version)
	assert.Equal(t, major, 0, "Invalid major part")
	assert.Equal(t, minor, 1, "Invalid major part")
	assert.Equal(t, patch, 0, "Invalid major part")
	assert.Equal(t, buildNumber, 0, "Invalid major part")

	version = "0.1.0-1"
	major, minor, patch, buildNumber = parseTag(version)
	assert.Equal(t, major, 0, "Invalid major part")
	assert.Equal(t, minor, 1, "Invalid major part")
	assert.Equal(t, patch, 0, "Invalid major part")
	assert.Equal(t, buildNumber, 1, "Invalid major part")
}
