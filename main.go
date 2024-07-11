package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const timeShift = 1 * time.Millisecond

func main() {
	file, err := os.Open("luffy.srt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanner.Bytes()
		fmt.Println(scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
