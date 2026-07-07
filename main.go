package main

import (
	"os"
	"runtime/debug"

	"github.com/kunchenguid/treehouse/v3/cmd"
)

var version = ""

func init() {
	if version == "" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		} else {
			version = "dev"
		}
	}
}

func main() {
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
