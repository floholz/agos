package main

import (
	"github.com/floholz/agos/internal/cli"
	"github.com/floholz/agos/internal/core"
)

func main() {
	agos := core.NewAgosApp()
	cli.LaunchCLI(agos)
}
