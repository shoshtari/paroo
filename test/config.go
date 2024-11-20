package test

import (
	"os"
	"path"

	"github.com/shoshtari/paroo/internal/configs"
)

func GetTestConfig() configs.ParooConfig {
	pwd := os.Getenv("PWD")
	currentPath := pwd
	running := true
	for running {
		files, err := os.ReadDir(currentPath)
		if err != nil {
			panic(err)
		}
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
	if err != nil {
		panic(err)
	}
	return res
}
