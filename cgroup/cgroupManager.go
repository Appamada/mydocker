package cgroup

import "github.com/Appamada/mydocker/cgroup/subsystem"

type CgroupManager struct {
	// cgroup在hierarchy中的路径
	Path string
	// 资源配置
	Resource *subsystem.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	manager := &CgroupManager{
		Path: path,
	}
	return manager
}

func (c *CgroupManager) Apply(pid int) error {
	for _, sysSysIns := range subsystem.SubsystemsIns {
		if err := sysSysIns.Apply(c.Path, pid); err != nil {
			return err
		}
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystem.ResourceConfig) error {
	for _, sysSysIns := range subsystem.SubsystemsIns {
		if err := sysSysIns.Set(c.Path, res); err != nil {
			return err
		}
	}
	return nil
}

func (c *CgroupManager) Destory() error {
	for _, sysSysIns := range subsystem.SubsystemsIns {
		if err := sysSysIns.Remove(c.Path); err != nil {
			return err
		}
	}
	return nil
}
