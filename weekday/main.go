package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	weekday := t.Weekday()

	dayDiff := int(weekday - 1)
	// 周日的情况
	if dayDiff < 0 {
		dayDiff = 6
	}

	startTime := t.Add(time.Duration(- dayDiff) * 24 * time.Hour)
	endTime := t.Add(time.Duration(4 - dayDiff) * 24 * time.Hour)

	fmt.Printf("%s - %s\n", startTime.Format("2006/01/02"), endTime.Format("2006/01/02"))
}
