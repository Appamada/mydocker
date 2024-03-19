package main

import (
	"fmt"

	"github.com/Appamada/mydocker/cgroups/subsystem"
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

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit container into image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}

		imageName := context.Args().Get(0)
		container.CommitImage(imageName)
		return nil
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
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container process",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "set volume, like: -v /tmp",
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
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
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

		detach := context.Bool("d")
		tty := context.Bool("ti")

		if detach && tty {
			return fmt.Errorf("it is not allowed to use tty and detach at the same time")
		}

		volume := context.String("v")
		containerName := context.String("name")

		resConf := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare:    context.String("cpu"),
			CpuSet:      context.String("cpuset"),
		}

		Run(containerName, tty, cmdArray, volume, resConf)
		return nil
	},
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list the info about containers",
	Action: func(context *cli.Context) error {
		container.ContainerRecordList()
		return nil
	},
}

var logCommand = cli.Command{
	Name:  "logs",
	Usage: "print logs of a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}

		containerName := context.Args().Get(0)
		container.LogContainer(containerName)
		return nil
	},
}
