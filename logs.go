package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	files := []string{"1","2","3","4","5"}
	lines := []string{"a","b","c","d","e"}

	for _, fileName := range files {
		filePath := "test" + string(os.PathSeparator) + fileName + ".log"

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Panic(err)
		}

		for _, lineName := range lines {
			currentTime := time.Now()
			lineText := fmt.Sprintf("%s %s %s\n", currentTime.Format(time.RFC3339), fileName, lineName)

			_, err = file.Write([]byte(lineText))
			if err != nil {
				log.Panic(err)
			}
		}

		err = file.Close()
		if err != nil {
			log.Panic(err)
		}

	}

}