package commands

import "os/exec"

func Exist(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}
