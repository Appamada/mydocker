package subsystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/Appamada/mydocker/util"
)

type MemorySubSystem struct {
}

func (m *MemorySubSystem) Apply(cgroupPath string, Pid int) error {
	if subsysCgroupPath, err := util.GetCgroupPath(m.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(Pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", subsysCgroupPath, err)
	}
}

func (m *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := util.GetCgroupPath(m.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			if err := ioutil.WriteFile(
				path.Join(subsysCgroupPath, "memory.limit_in_bytes"),
				[]byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", subsysCgroupPath, err)
	}
}

func (m *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := util.GetCgroupPath(m.Name(), cgroupPath, false); err == nil {
		if err := os.Remove(subsysCgroupPath); err != nil {
			return fmt.Errorf("remove cgroup %s error: %v", subsysCgroupPath, err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

func (m *MemorySubSystem) Name() string {
	return "memory"
}
