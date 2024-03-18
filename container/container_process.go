package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/Appamada/mydocker/util"
	log "github.com/sirupsen/logrus"
)

// const containerRootDir = "/root/busybox"

var (
	MntURL  = "/root/mnt/"
	RootURL = "/root/"
)

func NerParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe() //0，1，err
	if err != nil {
		log.Errorf("new pipe error %v", err)
		return nil, nil
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.ExtraFiles = []*os.File{readPipe} //将管道一端传入到容器进程中,容器进程接收数据

	NewWorkSapce(RootURL, MntURL)
	cmd.Dir = MntURL //设置工作目录且工作目录为空，导致无法找到/proc/self/exe。增加一个镜像tar包解决

	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func NewWorkSapce(rootURL string, mntURL string) {
	CreateReadLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountLayer(rootURL, mntURL)
}

func CreateReadLayer(rootURL string) {
	busyBoxUrl := rootURL + "busybox/"
	busyboxTarUrl := rootURL + "busybox.tar"

	exist, err := util.PathExists(busyBoxUrl)

	if err != nil {
		log.Errorf("check dir %s error %v", busyBoxUrl, err)
	}

	if !exist {
		if err := os.Mkdir(busyBoxUrl, 0777); err != nil {
			log.Errorf("create dir %s error %v", busyBoxUrl, err)
		}

	} else {
		log.Infof("busybox dir %s already exists", busyBoxUrl)
	}

	if _, err := exec.Command(
		"tar", "xvf", busyboxTarUrl, "-C", busyBoxUrl).CombinedOutput(); err != nil {
		log.Errorf("untar busybox tar file  %s error %v", busyboxTarUrl, err)
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"

	exist, err := util.PathExists(writeURL)

	if err != nil {
		log.Errorf("check dir %s error %v", writeURL, err)
	}

	if !exist {
		if err := os.Mkdir(writeURL, 0777); err != nil {
			log.Errorf("create dir %s error %v", writeURL, err)
		}
	}
}

func CreateMountLayer(rootURL string, mntURL string) {
	exist, err := util.PathExists(mntURL)

	if err != nil {
		log.Errorf("check dir %s error %v", mntURL, err)
	}

	if !exist {
		if err := os.Mkdir(mntURL, 0777); err != nil {
			log.Errorf("create dir %s error %v", mntURL, err)
		}
	}

	//mount -t aufs -o dirs=/root/writeLayer:/root/busybox none /root/mnt
	dirs := "dirs=" + rootURL + "writeLayer" + ":" + rootURL + "busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if cmd.Run(); err != nil {
		log.Errorf("mount %s error %v", mntURL, err)
	}

}

func DeleteWorkSpace(rootURL string, mntURL string) {
	UmountMntPoint(mntURL)
	DeleteWriteLayer(rootURL)
}

func UmountMntPoint(mntURL string) {
	cmd := exec.Command("umount", mntURL)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("umount %s error %v", mntURL, err)
	}

	if err := os.Remove(mntURL); err != nil {
		log.Errorf("remove %s error %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	if err := os.RemoveAll(rootURL + "writeLayer"); err != nil {
		log.Errorf("remove %s error %v", rootURL+"writeLayer", err)
	}
}
