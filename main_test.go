package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionTag_String(t *testing.T) {
	v123 := &VersionTag{
		ref:   nil,
		Tag:   "v",
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre:   "",
		Build: "",
	}

	assert.Equal(t, "v1.2.3", v123.String())

	v123p1 := &VersionTag{
		ref:   nil,
		Tag:   "v",
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre:   "pre.1",
		Build: "",
	}

	assert.Equal(t, "v1.2.3-pre.1", v123p1.String())

	v123p1b1 := &VersionTag{
		ref:   nil,
		Tag:   "v",
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre:   "pre.1",
		Build: "b1",
	}

	assert.Equal(t, "v1.2.3-pre.1+b1", v123p1b1.String())

	v123b1 := &VersionTag{
		ref:   nil,
		Tag:   "v",
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre:   "",
		Build: "b1",
	}

	assert.Equal(t, "v1.2.3+b1", v123b1.String())

	n123 := &VersionTag{
		ref:   nil,
		Tag:   "",
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre:   "",
		Build: "",
	}

	assert.Equal(t, "1.2.3", n123.String())
}

func TestVersionFromString(t *testing.T) {
	n123 := "1.2.3"
	v123 := "v1.2.3"
	v123p1 := "v1.2.3-p1"
	v123p1b1 := "v1.2.3-p1+b1"
	v123b1 := "v1.2.3+b1"

	v_n123, err := VersionFromString(n123)
	assert.NoError(t, err)
	assert.Equal(t, n123, v_n123.String())

	v_v123, err := VersionFromString(v123)
	assert.NoError(t, err)
	assert.Equal(t, v123, v_v123.String())

	v_v123p1, err := VersionFromString(v123p1)
	assert.NoError(t, err)
	assert.Equal(t, v123p1, v_v123p1.String())

	v_v123p1b1, err := VersionFromString(v123p1b1)
	assert.NoError(t, err)
	assert.Equal(t, v123p1b1, v_v123p1b1.String())

	v_v123b1, err := VersionFromString(v123b1)
	assert.NoError(t, err)
	assert.Equal(t, v123b1, v_v123b1.String())
}
