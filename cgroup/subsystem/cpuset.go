package subsystem

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/Appamada/mydocker/util"
)

type CpuSetSubSystem struct {
}

func (c *CpuSetSubSystem) Apply(cgroupPath string, Pid int) error {
	subSystenCgoupPath, err := util.GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}

	if err := os.WriteFile(path.Join(subSystenCgoupPath, "tasks"), []byte(strconv.Itoa(Pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc %d error: %v", Pid, err)
	}

	return nil
}

func (c *CpuSetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subSystenCgoupPath, err := util.GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}

	if err := os.WriteFile(path.Join(subSystenCgoupPath, "cpuset.cpus"), []byte(res.CpuShare), 0644); err != nil {
		return fmt.Errorf("set cgroup cpuset.cpus error: %v", err)
	}
	return nil
}

func (c *CpuSetSubSystem) Remove(cgroupPath string) error {
	subSystenCgoupPath, err := util.GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}

	if err := os.RemoveAll(subSystenCgoupPath); err != nil {
		return fmt.Errorf("remove cgroup %s error: %v", cgroupPath, err)
	}

	return nil
}

func (c *CpuSetSubSystem) Name() string {
	return "cpuset"
}
