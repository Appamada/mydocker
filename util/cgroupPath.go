package util

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

func FindCgroupRootPath(subSystemName string) (string, error) {
	f, err := os.Open("/proc/self/mountinfo")

	if err != nil {
		return "", err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Split(line, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			fmt.Println(opt)
			if opt == subSystemName {
				return fields[4], nil
			}
		}
	}

	if scanner.Err() != nil {
		return "", err
	}

	return "", err
}

func GetCgroupPath(subSystemName string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot, err := FindCgroupRootPath(subSystemName)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err != nil {
		if autoCreate {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
				return "", fmt.Errorf("create cgroup path %s error: %v", cgroupPath, err)
			}
			return path.Join(cgroupRoot, cgroupPath), nil
		} else {
			return "", fmt.Errorf("cgroup %s does not exist", cgroupPath)
		}

	}
	return path.Join(cgroupRoot, cgroupPath), nil
}
