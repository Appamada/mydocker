package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	cmdArry := readUserCommand()
	if cmdArry == nil || len(cmdArry) == 0 {
		return fmt.Errorf("run container get user command error")
	}

	path, err := exec.LookPath(cmdArry[0])
	if err != nil {
		log.Errorf("exec look path error %v", err)
		return err
	}

	log.Infof("find path %s", path)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// func Mount(source string, target string, fstype string, flags uintptr, data string) (err error)
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	if err := syscall.Exec(path, cmdArry[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")

	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Errorf("read pipe error %v", err)
		return nil
	}

	msgStr := string(msg)

	return strings.Split(msgStr, " ")
}
