package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/Appamada/mydocker/cgroup"
	"github.com/Appamada/mydocker/cgroup/subsystem"
	"github.com/Appamada/mydocker/container"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `mydocker is a simple container runtime implementation.
The purpose of this project is to learn how docker works and how to write a docker by ourselves
Enjoy it, just for fun.`

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
	mydocker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpu",
			Usage: "cpu limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}

		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		tty := context.Bool("ti")

		resConf := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare:    context.String("cpu"),
			CpuSet:      context.String("cpuset"),
		}

		Run(tty, cmdArray, resConf)
		return nil
	},
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("cmd is %s", cmdStr)
	writePipe.WriteString(cmdStr)
	writePipe.Close()
}

func Run(tty bool, cmdArray []string, resConfig *subsystem.ResourceConfig) {
	parent, writePipe := container.NerParentProcess(tty)

	if parent == nil {
		log.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Errorf("parent start error %v", err)
	}

	cgroupManager := cgroup.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destory()
	cgroupManager.Apply(parent.Process.Pid)
	cgroupManager.Set(resConfig)

	sendInitCommand(cmdArray, writePipe)
	parent.Wait()

	if err := syscall.Unmount("/proc", 0); err != nil {
		log.Error(err)
	}

	log.Infof("container process exit")
	os.Exit(0)

}

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(c *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
