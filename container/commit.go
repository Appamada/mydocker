package container

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func CommitImage(imageName string) {
	imageTar := "/root/" + imageName + ".tar"
	fmt.Println("tar file named ", imageTar)
	if _, err := exec.Command("tar", "czvf", imageTar, "-C", MntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("tar error %v", err)
	}
}
