package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	cmdArry := readUserCommand()
	if cmdArry == nil || len(cmdArry) == 0 {
		return fmt.Errorf("run container get user command error")
	}

	setMount()

	path, err := exec.LookPath(cmdArry[0])
	if err != nil {
		log.Errorf("exec look path error %v", err)
		return err
	}

	log.Infof("find path %s", path)

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

func setMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("get current work dir error %v", err)
		return
	}

	log.Infof("current work dir is %s", pwd)
	pivotRoot(pwd)

	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// func Mount(source string, target string, fstype string, flags uintptr, data string) (err error)
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir .pivot_root error: %v", err)
	}

	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root error: %v", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / error: %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")

	//umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount .pivot_root error: %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
