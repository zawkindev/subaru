package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	u "subaru/utility"
)

const timeShift int64 = 1
const filename = "luffy.srt"

func main() {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var newContent string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		yes, index := u.Includes("-->", line)
		if yes {
			timestamp1 := line[0 : index-1]
			timestamp2 := line[index+4 : len(line)-1]

			timestamp1 = u.TimeShift(timestamp1, timeShift)
			timestamp2 = u.TimeShift(timestamp2, timeShift)

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
