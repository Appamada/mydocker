package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

const ENV_EXEC_PID = "container_pid"
const ENV_EXEC_CMD = "container_command"

func GetContainerPid(containerName string) (string, error) {
	configDirPath := fmt.Sprintf(DefaultContainerRootPath + "/" + containerName)
	configFilePath := configDirPath + "/" + DefaultConfigName

	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Errorf("read file %s error %v", configFilePath, err)
		return "", err
	}

	var containerInfo containerInfo

	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		log.Errorf("unmarshal file %s error %v", configFilePath, err)
		return "", err
	}

	return containerInfo.Pid, nil

}

func ExecContainer(containerName string, cmdSlice []string) {
	containerPid, err := GetContainerPid(containerName)

	if err != nil {
		log.Errorf("get container %s pid error %v", containerName, err)
		return
	}
	log.Infof("pid is %s", containerPid)

	cmdStr := strings.Join(cmdSlice, " ")
	log.Infof("cmd is %s", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, containerPid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)

	if err := cmd.Run(); err != nil {
		log.Errorf("exec container %s error %v", containerName, err)
	}
}
