package subsystem

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/Appamada/mydocker/util"
)

type CpuSubSystem struct {
}

func (c *CpuSubSystem) Apply(cgorupPath string, Pid int) error {
	subSystenCgoupPath, err := util.GetCgroupPath(c.Name(), cgorupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgorupPath, err)
	}

	if err := os.WriteFile(path.Join(subSystenCgoupPath, "tasks"), []byte(strconv.Itoa(Pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc %d error: %v", Pid, err)
	}

	return nil
}

func (c *CpuSubSystem) Set(cgorupPath string, res *ResourceConfig) error {
	subSystenCgoupPath, err := util.GetCgroupPath(c.Name(), cgorupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgorupPath, err)
	}

	if err := os.WriteFile(path.Join(subSystenCgoupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
		return fmt.Errorf("set cgroup cpu.shares error: %v", err)
	}
	return nil
}

func (c *CpuSubSystem) Remove(cgorupPath string) error {
	subSystenCgoupPath, err := util.GetCgroupPath(c.Name(), cgorupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgorupPath, err)
	}

	if err := os.RemoveAll(subSystenCgoupPath); err != nil {
		return fmt.Errorf("remove cgroup %s error: %v", cgorupPath, err)
	}

	return nil
}

func (c *CpuSubSystem) Name() string {
	return "cpu"
}
