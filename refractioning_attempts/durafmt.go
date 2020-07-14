package main

import (
	"fmt"

	"github.com/hako/durafmt"
)

func main() {
	duration, err := durafmt.ParseString("354h22m3.24s")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(duration) // 2 weeks 18 hours 22 minutes 3 seconds
	// duration.String() // String representation. "2 weeks 18 hours 22 minutes 3 seconds"
}
