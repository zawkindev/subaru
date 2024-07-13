package main

import (
	u "subaru/utility"
)

func main() {
	timeShift := 1
	filename := "luffy.srt"

	u.TimeShift(filename, int64(timeShift))
}
