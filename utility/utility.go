package utility

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func timeStampToMilliseconds(timestamp string) (int64, error) {
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

func millisecondsToTimeStamp(milliseconds int64) string {
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

func includes(substr, str string) (bool, int) {
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

func timeAdd(timestamp string, millisecond int64) string {
	time, err := timeStampToMilliseconds(timestamp)
	if err != nil {
		log.Fatal(err)
	}
	return millisecondsToTimeStamp(time + millisecond)
}

func TimeShift(filename string, timeShift int64) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var newContent string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		yes, index := includes("-->", line)
		if yes {
			timestamp1 := line[0 : index-1]
			timestamp2 := line[index+4:len(line)-1]

			timestamp1 = timeAdd(timestamp1, timeShift)
			timestamp2 = timeAdd(timestamp2, timeShift)

			line = fmt.Sprintf("%s --> %s", timestamp1, timestamp2)
		}
		newContent += line + "\n"
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	file.Close()

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(newContent)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("file rewrited successfully")
}
