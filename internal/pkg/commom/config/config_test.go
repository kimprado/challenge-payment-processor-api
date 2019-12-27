// +build test unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCreateRedisDB(t *testing.T) {
	r := NewRedisDB(Configuration{})
	assert.NotNil(t, r)
	return
}

func TestNewCreateLoggingLevels(t *testing.T) {
	ll := NewLoggingLevels(Configuration{})
	assert.Nil(t, ll)
	return
}
