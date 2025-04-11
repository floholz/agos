package main

import (
	"github.com/floholz/agos/internal/core"
	"github.com/floholz/agos/internal/ui"
)

func main() {
	agos := core.NewAgosApp()
	ui.LaunchUI(agos)
}
