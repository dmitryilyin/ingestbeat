package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {

	fileHandle, _ := os.Open("1.txt")
	defer fileHandle.Close()

	reader := bufio.NewReader(fileHandle)
	batchMaxSize := 5

	var message strings.Builder
	var header strings.Builder
	var line string
	var err  error
	var batchSize = 0
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		if line == "\n" || line == "\r\n" {
			continue
		}
		if headerLine(&line, &header) {
			continue
		}
		if skipLine(&line) {
			continue
		}

		message.WriteString(line)

		if batchMaxSize > 0 {
			batchSize += 1
			if batchSize >= batchMaxSize {
				out(header.String(), message.String())
				message.Reset()
				batchSize = 0
			}
		}

	}
	out(header.String(), message.String())
}

func out(header string, message string)  {
	fmt.Println("= start =")
	fmt.Print(header)
	fmt.Print(message)
	fmt.Println("= end =")
}

func headerLine(line *string, header *strings.Builder) bool {
	headerPatterns := []string{`#\s*Header:`,`#\s*Software:`}
	for _, pattern := range headerPatterns {
		match, err := regexp.Match(pattern, []byte(*line))
		if err != nil {
			continue
		}
		if match {
			header.WriteString(*line)
			return true
		}
	}
	return false
}

func skipLine(line *string) bool {
	skipPatterns := []string{`^\s*#`}
	for _, pattern := range skipPatterns {
		match, err := regexp.Match(pattern, []byte(*line))
		if err != nil {
			continue
		}
		if match {
			return true
		}
	}
	return false
}
