package cgroups

import (
	"github.com/Appamada/mydocker/cgroups/subsystem"
	log "github.com/sirupsen/logrus"
)

type CgroupManager struct {
	// cgroup在hierarchy中的路径
	Path string
	// 资源配置
	Resource *subsystem.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
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
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
