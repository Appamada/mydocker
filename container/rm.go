package container

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func RmContainer(containerName string) error {
	configDirPath := fmt.Sprintf(DefaultContainerRootPath + "/" + containerName)

	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("get container %s info error %v", containerName, err)
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
