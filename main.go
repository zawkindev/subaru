package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	u "subaru/utility"
)

const timeShift = 1
const filename = "luffy.srt"

func main() {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		yes, index := u.Includes("-->", line)
		if yes {
			timestamp := line[0 : index-1]
			oldTimeInMilli, err := u.TimeStampToMilliseconds(timestamp)
			if err!=nil{
				log.Fatal(err)
			}
			shiftedTime := u.TimeShift(int64(timeShift), timestamp)
			fmt.Printf("%v --> %v \n", u.MillisecondsToTimeStamp(oldTimeInMilli), shiftedTime)
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
