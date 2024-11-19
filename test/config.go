package test

import (
	"os"
	"path"
	"testing"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/stretchr/testify/assert"
)

func GetTestConfig(t *testing.T) configs.ParooConfig {
	pwd := os.Getenv("PWD")
	currentPath := pwd
	running := true
	for running {
		files, err := os.ReadDir(currentPath)
		assert.Nil(t, err)
		for _, file := range files {
			if file.Name() == "go.mod" {
				running = false
				break
			}
		}

		if running {
			currentPath = path.Dir(currentPath)
		}
	}
	currentPath = path.Join(currentPath, "test")
	res, err := configs.GetConfig(currentPath)
	assert.Nil(t, err)
	return res
}
