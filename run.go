package main

import (
	"os"
	"strings"

	"github.com/Appamada/mydocker/cgroups"
	"github.com/Appamada/mydocker/cgroups/subsystem"
	"github.com/Appamada/mydocker/container"
	"github.com/Appamada/mydocker/util"
	log "github.com/sirupsen/logrus"
)

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("cmd is %s", cmdStr)
	writePipe.WriteString(cmdStr)
	writePipe.Close()
}

func Run(containerName string, tty bool, cmdArray []string, volume string, resConfig *subsystem.ResourceConfig, envArray []string) {
	var name string
	id := util.RandomString(10)

	if name != "" {
		name = containerName
	} else {
		name = id
	}

	parent, writePipe := container.NerParentProcess(tty, volume, &name, envArray)

	if parent == nil {
		log.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Errorf("parent start error %v", err)
	}

	container.ContainerInfoRecord(name, &id, parent.Process.Pid, cmdArray)

	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destory()

	if err := cgroupManager.Set(resConfig); err != nil {
		log.Errorf("set cgroup error %v", err)
	}

	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(cmdArray, writePipe)

	if tty {
		if err := parent.Wait(); err != nil {
			log.Errorf("parent wait error %v", err)
		}

		container.DeleteWorkSpace(container.RootURL, container.MntURL, volume)
		container.ContainerInfoDelete(name)
	}

	// if err := syscall.Unmount("/proc", 0); err != nil {
	// 	log.Error(err)
	// }

	// log.Infof("container process exit")
	// os.Exit(0)

}
