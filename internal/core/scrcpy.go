package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ScrCpy struct {
	ScrcpyCmd string
}

func NewScrCpy(scrcpyCmd string) *ScrCpy {
	scrcpy := &ScrCpy{
		ScrcpyCmd: scrcpyCmd,
	}
	return scrcpy
}

func (sc *ScrCpy) Exec(args ...string) error {
	execCmd := exec.Command(sc.ScrcpyCmd, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	return execCmd.Run()
}

func (sc *ScrCpy) ExecCallback(callback func(outBuf, errBuf strings.Builder) error, args ...string) error {
	execCmd := exec.Command(sc.ScrcpyCmd, args...)

	var outBuf, errBuf strings.Builder
	execCmd.Stdout = &outBuf
	execCmd.Stderr = &errBuf

	err := execCmd.Run()
	if err != nil {
		return err
	}
	return callback(outBuf, errBuf)
}

func (sc *ScrCpy) Run(ip string) error {
	err := sc.Exec("-s", ip)
	if err != nil {
		fmt.Println("Failed to run scrcpy:", err)
		return err
	}
	fmt.Println("scrcpy session ended.")
	return nil
}
