package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadSpecificEnvFile(envFile string) {
	goModDirPath := findGoModDir(searchCallerFile())

	isLocal := goModDirPath != ""

	if isLocal {
		envFile = filepath.Join(findGoModDir(searchCallerFile()), envFile)
	}

	err := godotenv.Overload(envFile)
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	if err != nil {
		panic(fmt.Sprintf("Задан env файл %s, но его не удалось загрузить: %s", envFile, err.Error()))
	}
}

func LoadEnv() {
	LoadSpecificEnvFile(".env")
	LoadSpecificEnvFile(".env.override")
}

func searchCallerFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	currentFile := file
	for i := 2; file == currentFile; i++ {
		_, file, _, ok = runtime.Caller(i)
		if !ok {
			return ""
		}
	}

	return file
}

func findGoModDir(from string) string {
	dir := filepath.Dir(from)
	gopath := filepath.Clean(os.Getenv("GOPATH"))
	for filepath.Dir(dir) != dir && dir != gopath {
		goModFile := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModFile); os.IsNotExist(err) {
			dir = filepath.Dir(dir)
			continue
		}
		return dir
	}
	return ""
}
