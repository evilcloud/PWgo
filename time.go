package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

func humaniseDuration(start time.Time) string {
	ret := humanize.Time(start)
	debugNotification(ret)
	if ret == "a long while ago" {
		ret = ""
	}
	return ret
}

func executableLastMod() string {
	filename, err := os.Executable()
	if err != nil {
		return "Failed to capture executable info"
	}
	file, err := os.Stat(filename)
	if err != nil {
		return "Failed to capture mod time"
	}
	return fmt.Sprint(file.ModTime())
}
