package util

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func FileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateFolder(path string) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		logrus.Errorf("Create folder %s error: %v", path, err)
	}
	return err
}

func CreateFile(path string) (*os.File, error) {
	basePath := filepath.Dir(path)
	if err := CreateFolder(basePath); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func JsonToFile(dst string, data interface{}) error {
	str, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logrus.Errorf("JsonToFile error: %v", err)
		return err
	}
	err = os.WriteFile(dst, str, 0700)
	if err != nil {
		logrus.Errorf("WriteFile error: %v", err)
		return err
	}
	return nil
}
