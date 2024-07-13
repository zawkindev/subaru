package utility

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func TimeStampToMilliseconds(timestamp string) (int64, error) {
	parts := strings.Split(timestamp, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid timestamp format")
	}

	secondsPart := strings.Split(parts[2], ",")
	if len(secondsPart) != 2 {
		return 0, fmt.Errorf("invalid timestamp format")
	}

	arrStr := [4]string{}
	arrInt := [4]int64{}
	var err error

	arrStr[0], arrStr[1], arrStr[2], arrStr[3] = parts[0], parts[1], secondsPart[0], secondsPart[1]
	/*
	   0 -> hours
	   1 -> minutes
	   2 -> seconds
	   3 -> milliseconds
	*/

	for i, v := range arrStr {
		arrInt[i], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("converting string to int failed")
		}
	}

	return arrInt[0]*time.Hour.Milliseconds() + arrInt[1]*time.Minute.Milliseconds() + arrInt[2]*time.Second.Milliseconds() + arrInt[3], nil
}

func MillisecondsToTimeStamp(milliseconds int64) string {
	arrInt := [4]int64{}
	/*
	   0 -> hours
	   1 -> minutes
	   2 -> seconds
	   3 -> milliseconds
	*/

	arrInt[0] = milliseconds / time.Hour.Milliseconds()
	milliseconds -= arrInt[0] * time.Hour.Milliseconds()

	arrInt[1] = milliseconds / time.Minute.Milliseconds()
	milliseconds -= arrInt[1] * time.Minute.Milliseconds()

	arrInt[2] = milliseconds / time.Second.Milliseconds()
	milliseconds -= arrInt[2] * time.Second.Milliseconds()

	arrInt[3] = milliseconds

	timestamp := fmt.Sprintf("%02d:%02d:%02d,%03d", arrInt[0], arrInt[1], arrInt[2], arrInt[3])

	return timestamp
}

func Includes(substr, str string) (bool, int) {
	start := 0
	end := len(substr)
	for end < len(str) {
		if str[start:end] == substr {
			return true, start
		}
		start++
		end++
	}
	return false, 0
}

func TimeShift(timestamp string, millisecond int64) string {
	time, err := TimeStampToMilliseconds(timestamp)
	if err != nil {
		log.Fatal(err)
	}
	return MillisecondsToTimeStamp(time + millisecond)
}
