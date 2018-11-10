package db

import (
	"github.com/marvin-automator/marvin/internal/config"
	"os"
)

func SetupTestDB() {
	config.DataDir = "./test_data"
	os.MkdirAll(config.DataDir, os.ModePerm)
}

func TearDownTestDB() {
	os.RemoveAll(config.DataDir)
}
