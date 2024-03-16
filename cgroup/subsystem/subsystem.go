package subsystem

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type SubSystem interface {
	Apply(cgorupPath string, pid int) error
	Set(cgorupPath string, res *ResourceConfig) error
	Remove(cgorupPath string) error
	Name() string
}

var (
	SubsystemsIns = []SubSystem{
		&CpuSubSystem{},
		&CpuSetSubSystem{},
		&MemorySubSystem{},
	}
)
