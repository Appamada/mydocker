package main

import (
	"os"
	"strings"

	"github.com/Appamada/mydocker/cgroups"
	"github.com/Appamada/mydocker/cgroups/subsystem"
	"github.com/Appamada/mydocker/container"
	log "github.com/sirupsen/logrus"
)

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("cmd is %s", cmdStr)
	writePipe.WriteString(cmdStr)
	writePipe.Close()
}

func Run(tty bool, cmdArray []string, resConfig *subsystem.ResourceConfig) {
	parent, writePipe := container.NerParentProcess(tty)

	if parent == nil {
		log.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Errorf("parent start error %v", err)
	}

	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destory()

	if err := cgroupManager.Set(resConfig); err != nil {
		log.Errorf("set cgroup error %v", err)
	}

	if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
		log.Errorf("apply cgroup error %v", err)
	}

	sendInitCommand(cmdArray, writePipe)

	if tty {
		if err := parent.Wait(); err != nil {
			log.Errorf("parent wait error %v", err)
		}

		container.DeleteWorkSpace(container.RootURL, container.MntURL)
	}

	// if err := syscall.Unmount("/proc", 0); err != nil {
	// 	log.Error(err)
	// }

	// log.Infof("container process exit")
	// os.Exit(0)

}
