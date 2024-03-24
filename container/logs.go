package container

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func LogContainer(containerName string) {
	containerLogPath := fmt.Sprintf(DefaultContainerRootPath + "/" + containerName + "/" + DefaultLogName)

	file, err := os.Open(containerLogPath)
	defer file.Close()

	if err != nil {
		log.Errorf("open file %s error %v", containerLogPath, err)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("read file %s error %v", containerLogPath, err)
		return
	}

	fmt.Fprint(os.Stdout, string(content))
}
