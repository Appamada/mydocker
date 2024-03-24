package main

import (
	"fmt"
	"os"

	"github.com/Appamada/mydocker/cgroups/subsystem"
	"github.com/Appamada/mydocker/container"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/Appamada/mydocker/nsenter"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a running container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}

		containerName := context.Args().Get(0)
		container.StopContainer(containerName)
		return nil
	},
}

var startCommand = cli.Command{
	Name:  "start",
	Usage: "start a stopped container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}

		containerName := context.Args().Get(0)
		container.StartContainer(containerName)
		return nil
	},
}

var rmCommand = cli.Command{
	Name:  "rm",
	Usage: "remove unused containers which is in stopped state",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}

		containerName := context.Args().Get(0)
		container.RmContainer(containerName)
		return nil
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

const ENV_EXEC_PID = "container_pid"
const ENV_EXEC_CMD = "container_command"

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "execute a command in container",
	// Flags: []cli.Flag{
	// 	cli.BoolFlag{
	// 		Name:  "ti",
	// 		Usage: "enable tty",
	// 	},
	// },
	Action: func(context *cli.Context) error {
		if os.Getenv(ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %s", os.Getgid())
			return nil
		}

		if len(context.Args()) < 2 {
			log.Errorf("missing container id or command")
			return nil
		}

		var cmdSlice []string
		for _, cmd := range context.Args().Tail() {
			cmdSlice = append(cmdSlice, cmd)
		}

		containerName := context.Args().Get(0)

		container.ExecContainer(containerName, cmdSlice)
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
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "set environment variables, like: -e LANG=zh_CN.UTF-8",
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

		envSlice := context.StringSlice("e")
		volume := context.String("v")
		containerName := context.String("name")

		resConf := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare:    context.String("cpu"),
			CpuSet:      context.String("cpuset"),
		}

		Run(containerName, tty, cmdArray, volume, resConf, envSlice)
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
