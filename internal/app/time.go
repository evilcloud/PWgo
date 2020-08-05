package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

func humaniseDuration(start time.Time) string {
	ret := humanize.Time(start)
	// debugNotification(ret)
	if ret == "a long while ago" {
		ret = ""
	}
	return ret
}

func executableLastMod() string {
	return fmt.Sprint(executableLastModTime())
}

func executableLastModTime() time.Time {
	filename, err := os.Executable()
	if err != nil {
		// popupMessage("Error", "Failed to capture executable info")
		log.Println("Failed to capture executable info")
	}
	file, err := os.Stat(filename)
	if err != nil {
		// popupMessage("Error", "Failed to capture mod time")
		log.Println("Failed to capture mod time")
	}
	return file.ModTime()
}
