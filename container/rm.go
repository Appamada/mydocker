package container

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func RmContainer(containerName string) error {
	configDirPath := fmt.Sprintf(DefaultContainerRootPath + "/" + containerName)
	configFilePath := configDirPath + "/" + DefaultConfigName

	byteContent, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Errorf("read file %s error %v", configFilePath, err)
		return err
	}

	var containerInfo containerInfo
	if err := json.Unmarshal(byteContent, &containerInfo); err != nil {
		log.Errorf("unmarshal file %s error %v", configFilePath, err)
		return err
	}

	if containerInfo.Status != STOP {
		log.Errorf("container %s is running", containerName)
		return fmt.Errorf("container %s is running", containerName)
	}

	if err := os.RemoveAll(configDirPath); err != nil {
		log.Errorf("remove dir %s error %v", configDirPath, err)
		return err
	}

	return nil
}
