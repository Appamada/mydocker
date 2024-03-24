package container

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func StartContainer(containerName string) error {

	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("get container %s info error %v", containerName, err)
		return err
	}

	if containerInfo.Status != STOP {
		log.Errorf("container %s is not in stopped status", containerName)
		return fmt.Errorf("container %s is not in stopped status", containerName)
	}

	cmd := exec.Command(containerInfo.Cmd)

	if err := cmd.Start(); err != nil {
		log.Errorf("start container %s error %v", containerName, err)
		return err
	}

	containerInfo.Status = RUNNING
	containerInfo.Pid = fmt.Sprintf("%d", cmd.Process.Pid)

	newContentByte, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("json marshal container info error %v", err)
		return err
	} else {
		configFilePath := fmt.Sprintf(DefaultContainerRootPath + "/" + containerName + "/" + DefaultConfigName)
		os.WriteFile(configFilePath, newContentByte, 0644)

		log.Infof("start container %s success", containerName)
		return nil
	}
}
