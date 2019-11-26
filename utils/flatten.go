package utils

import (
	"fmt"
	"os/exec"
)

func Flatten(folder string) error {
	t := fmt.Sprintf("%s/", folder)
	cmd := exec.Command("find", t, "-mindepth", "2", "-type", "f", "-exec", "mv", "-i", "{}", t, ";")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	RemoveEmptyFolder(folder)
	return nil
}
