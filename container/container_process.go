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
	MntURL  = "/root/mnt"
	RootURL = "/root"
)

func NerParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
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

	NewWorkSapce(RootURL, MntURL, volume)
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

func NewWorkSapce(rootURL string, mntURL string, volume string) {
	CreateReadLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountLayer(rootURL, mntURL)

	if volume != "" {
		volume := util.VolumeUrlExtract(volume)

		if len(volume) == 2 && volume[0] != "" && volume[1] != "" {
			MountVolume(volume, mntURL)
			log.Infof("mount volume %s to %s", volume[0], volume[1])
		} else {
			log.Infof("volume parameter is not correct")
		}
	}

}

func CreateReadLayer(rootURL string) {
	busyBoxUrl := rootURL + "/busybox/"
	busyboxTarUrl := rootURL + "/busybox.tar"

	exist, err := util.PathExists(busyBoxUrl)

	if err != nil {
		log.Errorf("check dir %s error %v", busyBoxUrl, err)
	}

	if !exist {
		util.PathCreate(busyBoxUrl)
	}

	if _, err := exec.Command(
		"tar", "xvf", busyboxTarUrl, "-C", busyBoxUrl).CombinedOutput(); err != nil {
		log.Errorf("untar busybox tar file  %s error %v", busyboxTarUrl, err)
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "/writeLayer"

	exist, err := util.PathExists(writeURL)

	if err != nil {
		log.Errorf("check dir %s error %v", writeURL, err)
	}

	if !exist {
		util.PathCreate(writeURL)
	}
}

func CreateMountLayer(rootURL string, mntURL string) {
	exist, err := util.PathExists(mntURL)

	if err != nil {
		log.Errorf("check dir %s error %v", mntURL, err)
	}

	if !exist {
		util.PathCreate(mntURL)
	}

	//mount -t aufs -o dirs=/root/writeLayer:/root/busybox none /root/mnt
	dirs := "dirs=" + rootURL + "/writeLayer" + ":" + rootURL + "/busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if cmd.Run(); err != nil {
		log.Errorf("mount %s error %v", mntURL, err)
	}

}

func MountVolume(volume []string, mntPoint string) {
	parentVolumeURL := volume[0]
	exist, err := util.PathExists(parentVolumeURL)

	if err != nil {
		log.Errorf("check dir %s error %v", parentVolumeURL, err)
	}

	if !exist {
		util.PathCreate(parentVolumeURL)
	}

	containerVolumeURL := mntPoint + volume[1]
	exist, err = util.PathExists(containerVolumeURL)

	if err != nil {
		log.Errorf("check dir %s error %v", containerVolumeURL, err)
	}

	if !exist {
		util.PathCreate(containerVolumeURL)
	}

	dirs := "dirs=" + parentVolumeURL
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("mount %s error %v", containerVolumeURL, err)
	}
}

func DeleteWorkSpace(rootURL string, mntURL string, volume string) {
	if volume != "" {
		volume := util.VolumeUrlExtract(volume)

		if len(volume) == 2 && volume[0] != "" && volume[1] != "" {
			UmountMntPointWithVolume(mntURL, volume)
		} else {
			log.Infof("volume parameter is not correct")
		}
	}

	UmountMntPoint(mntURL)
	DeleteWriteLayer(rootURL)
}

func UmountMntPointWithVolume(mntURL string, volume []string) {
	containerVolumePath := mntURL + volume[1]

	cmd := exec.Command("umount", containerVolumePath)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Infof("umount %s error %v", containerVolumePath, err)
	}

	if err := os.RemoveAll(containerVolumePath); err != nil {
		log.Infof("remove %s error %v", containerVolumePath, err)
	}
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
