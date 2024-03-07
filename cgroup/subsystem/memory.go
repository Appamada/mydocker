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

func (m *MemorySubSystem) Apply(cgorupPath string, Pid int) error {
	subSystenCgoupPath, err := util.GetCgroupPath(m.Name(), cgorupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgorupPath, err)
	}

	if err := os.WriteFile(path.Join(subSystenCgoupPath, "tasks"), []byte(strconv.Itoa(Pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc %d error: %v", Pid, err)
	}

	return nil
}

func (m *MemorySubSystem) Set(cgorupPath string, res *ResourceConfig) error {
	subSystenCgoupPath, err := util.GetCgroupPath(m.Name(), cgorupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgorupPath, err)
	}

	if err := os.WriteFile(path.Join(subSystenCgoupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
		return fmt.Errorf("set cgroup memory.limit_in_bytes error: %v", err)
	}
	return nil
}

func (m *MemorySubSystem) Remove(cgorupPath string) error {
	subSystenCgoupPath, err := util.GetCgroupPath(m.Name(), cgorupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgorupPath, err)
	}

	if err := os.RemoveAll(subSystenCgoupPath); err != nil {
		return fmt.Errorf("remove cgroup %s error: %v", cgorupPath, err)
	}

	return nil
}

func (m *MemorySubSystem) Name() string {
	return "memory"
}
