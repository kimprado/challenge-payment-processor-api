// +build test unit

package config

import (
	"testing"
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
