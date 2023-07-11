package commands

import (
	"testing"
)

func TestExist(t *testing.T) {
	if Exist("ls") == false {
		t.Error("ls command does not exist!")
	}
	if Exist("ls111") == true {
		t.Error("ls111 command should not exist!")
	}
}
