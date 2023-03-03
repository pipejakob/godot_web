package main

import (
	"fmt"
	"runtime/debug"
)

const VERSION = "0.0"

func version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return fmt.Sprintf("%s commit:%s", VERSION, setting.Value)
			}
		}
	}

	return VERSION
}
