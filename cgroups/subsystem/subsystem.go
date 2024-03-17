package subsystem

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type SubSystem interface {
	Apply(cgroupPath string, pid int) error
	Set(cgroupPath string, res *ResourceConfig) error
	Remove(cgroupPath string) error
	Name() string
}

var (
	SubsystemsIns = []SubSystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpuSetSubSystem{},
	}
)
