package container

import (
	"fmt"
	"io/ioutil"
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

	msg, err := ioutil.ReadAll(pipe)
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

	if err := pivotRoot(pwd); err != nil {
		log.Errorf("pivot root error %v", err)
		return
	}

	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// func Mount(source string, target string, fstype string, flags uintptr, data string) (err error)
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	// systemd 加入linux之后, mount namespace 就变成 shared by default, 必须显式声明新的mount namespace独立。
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return err
	}

	// 重新mount root
	// bind mount：将相同内容换挂载点
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir .pivot_root error: %v", err)
	}

	// pivot_root 到新的rootfs, 老的 old_root挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root error: %v", err)
	}

	// 修改当前的工作目录到根目录
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
