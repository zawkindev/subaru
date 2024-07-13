package main

import (
		u "subaru/utility"
)

const timeShift int64 = 1
const filename = "luffy.srt"

func main() {
	u.TimeShift(filename, timeShift)
}
