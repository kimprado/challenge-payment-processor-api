// +build test unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFlagsPadrao(t *testing.T) {
	expect := struct {
		configPath string
	}{
		configPath: "./configs/config-dev.json",
	}

	c := loadFlags()

	if c != expect.configPath {
		t.Errorf("Caminho esperado %q é diferente sde %q\n", expect.configPath, c)
	}

}

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
