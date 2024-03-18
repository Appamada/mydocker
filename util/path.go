package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathCreate(path string) {
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Errorf("create dir %s error %v", path, err)
	} else {
		log.Infof("parent dir %s already done", path)
	}
}
