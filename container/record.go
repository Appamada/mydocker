package container

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Appamada/mydocker/util"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultConfigRootPath = "/var/run/containerInfo"
	DefaultConfigName     = "config.json"

	RUNNING string = "running"
	STOP    string = "stopped"
	EXIT    string = "exit"
)

type containerInfo struct {
	Pid         string    `json:"pid"`
	Name        string    `json:"name"`
	CreatedTime time.Time `json:"start_time"`
	Cmd         string    `json:"cmd"`
	Status      string    `json:"status"`
}

func ContainerInfoDelete(name string) {
	configPath := fmt.Sprint(DefaultConfigRootPath + "/" + name)
	if err := os.RemoveAll(configPath); err != nil {
		log.Errorf("remove %s error %v", configPath, err)
	}
}

func ContainerInfoRecord(name string, pid int, cmdArray []string) {

	var containerName string

	command := strings.Join(cmdArray, "")

	id := util.RandomString(10)

	if name != "" {
		containerName = name
	} else {
		containerName = id
	}

	info := &containerInfo{
		Pid:         strconv.Itoa(pid),
		Name:        containerName,
		CreatedTime: time.Now(),
		Cmd:         command,
		Status:      RUNNING,
	}

	jsonBytes, err := json.Marshal(info)
	if err != nil {
		log.Errorf("record container info error %v", err)
	}

	jsonStr := string(jsonBytes)

	configPath := fmt.Sprintf(DefaultConfigRootPath + "/" + containerName)
	exist, err := util.PathExists(configPath)

	if err != nil {
		log.Errorf("check dir %s error %v", configPath, err)
	}

	if !exist {
		util.PathCreate(configPath)
	}

	configFileName := configPath + "/" + DefaultConfigName

	file, err := os.Create(configFileName)
	defer file.Close()

	if err != nil {
		log.Errorf("create file %s error %v", configFileName, err)
	}

	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("write file %s error %v", configFileName, err)
	}

}
