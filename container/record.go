package container

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Appamada/mydocker/util"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultContainerRootPath = "/var/run/containerInfo"
	DefaultConfigName        = "config.json"
	DefaultLogName           = "container.log"

	RUNNING string = "running"
	STOP    string = "stopped"
	EXIT    string = "exit"
)

type containerInfo struct {
	ID          string    `json:"id"`
	Pid         string    `json:"pid"`
	Name        string    `json:"name"`
	CreatedTime time.Time `json:"start_time"`
	Cmd         string    `json:"cmd"`
	Status      string    `json:"status"`
}

func getContainerInfoByName(containerName string) (*containerInfo, error) {

	configFilePath := fmt.Sprint(DefaultContainerRootPath + "/" + containerName + "/" + DefaultConfigName)

	content, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		log.Errorf("read file %s error %v", configFilePath, err)
		return nil, err
	}

	var containerInfo containerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("unmarshal file %s error %v", configFilePath, err)
		return nil, err
	}

	return &containerInfo, nil
}

func getContainerInfoByFile(file fs.FileInfo) (*containerInfo, error) {
	containerName := file.Name()

	configFilePath := fmt.Sprint(DefaultContainerRootPath + "/" + containerName + "/" + DefaultConfigName)

	content, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		log.Errorf("read file %s error %v", configFilePath, err)
		return nil, err
	}

	var containerInfo containerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("unmarshal file %s error %v", configFilePath, err)
		return nil, err
	}

	return &containerInfo, nil
}

func ContainerRecordList() {
	configPath := fmt.Sprint(DefaultContainerRootPath)

	files, err := ioutil.ReadDir(configPath)
	if err != nil {
		log.Errorf("read dir %s error %v", configPath, err)
		return
	}

	var containers []*containerInfo

	for _, file := range files {
		tmpContainer, err := getContainerInfoByFile(file)
		if err != nil {
			log.Errorf("get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	//使用tabwriter.NewWriter()在控制台打印出容器信息
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)

	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")

	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			item.Name,
			item.Pid,
			item.Status,
			item.Cmd,
			item.CreatedTime,
		)
	}

	if err := w.Flush(); err != nil {
		log.Errorf("flush error %v", err)
		return
	}
}

func ContainerInfoDelete(name string) {
	configPath := fmt.Sprint(DefaultContainerRootPath + "/" + name)
	if err := os.RemoveAll(configPath); err != nil {
		log.Errorf("remove %s error %v", configPath, err)
	}
}

func ContainerInfoRecord(name string, id *string, pid int, cmdArray []string) {

	command := strings.Join(cmdArray, "")

	info := &containerInfo{
		ID:          *id,
		Name:        name,
		Pid:         strconv.Itoa(pid),
		Status:      RUNNING,
		Cmd:         command,
		CreatedTime: time.Now(),
	}

	jsonBytes, err := json.Marshal(info)
	if err != nil {
		log.Errorf("record container info error %v", err)
	}

	jsonStr := string(jsonBytes)

	configPath := fmt.Sprintf(DefaultContainerRootPath + "/" + name)
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
