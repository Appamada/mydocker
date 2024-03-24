package container

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func StopContainer(containerName string) {
	pid, err := GetContainerPid(containerName)

	if err != nil {
		log.Errorf("get container %s pid error %v", containerName, err)
		return
	}

	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorf("convert pid %s to int error %v", pid, err)
		return
	}

	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("kill -15 %s error %v", pid, err)
		return
	}

	containerDir := fmt.Sprint(DefaultContainerRootPath + "/" + containerName)
	configFilePath := containerDir + "/" + DefaultConfigName

	byteContent, err := os.ReadFile(configFilePath)

	if err != nil {
		log.Errorf("read file %s error %v", configFilePath, err)
		return
	}

	var containerInfo containerInfo
	if err := json.Unmarshal(byteContent, &containerInfo); err != nil {
		log.Errorf("unmarshal file %s error %v", configFilePath, err)
		return
	}

	containerInfo.Status = STOP
	containerInfo.Pid = ""

	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("marshal file %s error %v", configFilePath, err)
		return
	}

	if err := os.WriteFile(configFilePath, newContentBytes, 0622); err != nil {
		log.Errorf("write file %s error %v", configFilePath, err)
		return
	}
}
