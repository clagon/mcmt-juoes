package config

import (
	"os"
)

func GetServerDir() string {
	if dir := os.Getenv("SERVER_DIR"); dir != "" {
		return dir
	}
	return "../server"
}

func GetDataDir() string {
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir
	}
	return "../data"
}
