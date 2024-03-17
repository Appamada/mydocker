package subsystem

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/Appamada/mydocker/util"
)

type MemorySubSystem struct {
}

func (m *MemorySubSystem) Apply(cgroupPath string, Pid int) error {
	subSystenCgoupPath, err := util.GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}

	if err := os.WriteFile(path.Join(subSystenCgoupPath, "tasks"), []byte(strconv.Itoa(Pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc %d error: %v", Pid, err)
	}

	return nil
}

func (m *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subSystenCgoupPath, err := util.GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}

	if err := os.WriteFile(path.Join(subSystenCgoupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
		return fmt.Errorf("set cgroup memory.limit_in_bytes error: %v", err)
	}
	return nil
}

func (m *MemorySubSystem) Remove(cgroupPath string) error {
	subSystenCgoupPath, err := util.GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}

	if err := os.RemoveAll(subSystenCgoupPath); err != nil {
		return fmt.Errorf("remove cgroup %s error: %v", cgroupPath, err)
	}

	return nil
}

func (m *MemorySubSystem) Name() string {
	return "memory"
}
