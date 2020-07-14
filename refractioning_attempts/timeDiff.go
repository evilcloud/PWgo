package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/hako/durafmt"
)

func main() {
	start := time.Now()
	end := start.Add(time.Second * 635342643)
	humaniseDuration(start, end)
	// diff := humaniseDuration(start, end)
	// fmt.Println(diff, reflect.TypeOf(diff))
}

func humaniseDuration(start, end time.Time) {
	diff := end.Sub(start)
	fmt.Println(diff, reflect.TypeOf(diff))

	duration, err := durafmt.ParseString(diff.String())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(duration)

	// days := diff.Round(time.Second)
	// hours := days / time.Hour
	// minutes := days / time.Minute
	// seconds := days / time.Second

	// fmt.Println("day", days)
	// fmt.Println("hours", hours)
	// fmt.Println("minute", minutes)
	// fmt.Println("second", seconds)
}
