package main

import (
	"os"
	"log"
	"path/filepath"
	"strings"
	"time"
	"fmt"
	"bufio"
	"io/ioutil"
)

func findCheckMarked(fileInfo os.FileInfo) bool {
	log.Println("FindCheckMarked: " + fileInfo.Name())
	if fileInfo.Size() == 0 {
		return false
	}
	return true
}

func findCheckMatch(fileInfo os.FileInfo, filePattern string) bool {
	log.Println("FindCheckMatch: " + fileInfo.Name())
	if fileInfo.IsDir() {
		return false
	}
	matched, err := filepath.Match(filePattern, fileInfo.Name())
	if err != nil {
		log.Println(err)
		return false
	}
	return matched
}

//func findFilesGlob(instance string, fileDirectory string, filePatterns []string, fileNamesChannel chan<- string) {
//	if instance == "" {
//		instance = "main"
//	}
//	log.Printf("FindFilesGlob[%s]: start", instance)
//
//	for _, pattern := range filePatterns {
//		matchedFiles, err := filepath.Glob(fileDirectory + string(os.PathSeparator) + pattern)
//		if err != nil {
//			log.Panic(err)
//		}
//
//		for _, matchedFile := range matchedFiles {
//			matchedFileInfo, err := os.Stat(matchedFile)
//			if err != nil {
//				log.Panic(err)
//			}
//
//			if !findCheckMarked(matchedFileInfo) {
//				continue
//			}
//
//			log.Printf("FindFilesGlob[%s]: send \"%s\"", instance, matchedFile)
//			fileNamesChannel <- matchedFile
//		}
//
//	}
//	log.Printf("FindFilesGlob[%s]: closing channel", instance)
//	close(fileNamesChannel)
//}

func findFilesReadDir(instance string, fileDirectory string, filePatterns []string, fileNamesChannel chan<- string) {
	if instance == "" {
		instance = "main"
	}
	log.Printf("FindFilesReadDir[%s]: start", instance)
	for {
		log.Printf("FindFilesReadDir[%s]: scan \"%s\"", instance, fileDirectory)
		fileInfoList, err := ioutil.ReadDir(fileDirectory)
		if err != nil {
			log.Panic(err)
		}
		for _, fileInfo := range fileInfoList {
			for _, filePattern := range filePatterns {

				if !findCheckMatch(fileInfo, filePattern) {
					continue
				}

				if !findCheckMarked(fileInfo) {
					continue
				}

				fileFullPath := fileDirectory + string(os.PathSeparator) + fileInfo.Name()
				log.Printf("FindFilesReadDir[%s]: send \"%s\"", instance, fileFullPath)
				fileNamesChannel <- fileFullPath
			}
		}

		time.Sleep(3 * time.Second)
	}
}

//func findFilesWalk(instance string, fileDirectory string, filePatterns []string, fileNamesChannel chan<- string) {
//	if instance == "" {
//		instance = "main"
//	}
//	log.Printf("FindFilesWalk[%s]: start", instance)
//
//}

func findTest(instance string, word string, fileNamesChannel chan<- string) {
	if instance == "" {
		instance = "main"
	}
	log.Printf("FindTest[%s]: start", instance)

	if word == "" {
		word = instance
	}

	var fileName string
	for index := 1; index <= 10; index++ {
		fileName = fmt.Sprintf("%s-%d", word, index)
		log.Printf("FindTest[%s]: send \"%s\"", instance, fileName)
		fileNamesChannel <- fileName
		time.Sleep(time.Second * 1)
	}
}

func processTest(instance string, fileNamesChannel <-chan string, eventsChannel chan<- string) {
	if instance == "" {
		instance = "main"
	}
	log.Printf("ProcessTest[%s]: start", instance)

	for fileName := range fileNamesChannel {
		log.Printf("ProcessTest[%s]: receive \"%s\"", instance, fileName)
		event := "Data from file: " + fileName
		log.Printf("ProcessTest[%s]: send \"%s\"", instance, event)
		eventsChannel <- event
	}
}

func processReadFile(instance string, fileNamesChannel <-chan string, eventsChannel chan<- string) {
	if instance == "" {
		instance = "main"
	}
	log.Printf("ProcessReadFile[%s]: start", instance)

	for filePath := range fileNamesChannel {
		log.Printf("ProcessReadFile[%s]: receive \"%s\"", instance, filePath)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("ProcessReadFile[%s]: send %d bytes", instance, len(content))
		eventsChannel <- string(content)
		processMarkDelete(filePath)
	}
}

func processMarkRename(filePath string) {
	newFilePath := filePath + ".processed"
	log.Printf("ProcessMarkRename: \"%s\" -> \"%s\"", filePath, newFilePath)
	err := os.Rename(filePath, newFilePath)
	if err != nil {
		log.Panic(err)
	}
}

func processMarkDelete(filePath string) {
	log.Printf("ProcessMarkDelete: \"%s\"", filePath)
	err := os.Remove(filePath)
	if err != nil {
		log.Panic(err)
	}
}

func outputTest(instance string, eventsChannel <-chan string) {
	if instance == "" {
		instance = "main"
	}
	log.Printf("OutputTest[%s]: start", instance)

	for event := range eventsChannel {
		log.Printf("OutputTest[%s]: receive \"%s\"", instance, event)
	}
}

func outputWriteFile(instance string, eventsChannel <-chan string) {
	if instance == "" {
		instance = "main"
	}
	log.Printf("OutputWriteFile[%s]: start", instance)

	file, err := os.OpenFile(instance + ".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}

	for event := range eventsChannel {
		if !strings.HasSuffix(event, "\n") {
			event += "\n"
		}
		_, err = file.Write([]byte(event))
		if err != nil {
			log.Panic(err)
		}
	}

	err = file.Close()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	filesDir := "test"
	filesPatterns := []string{"*.log"}

	var fileNamesChannel chan string = make(chan string, 0)
	var eventsChannel chan string = make(chan string, 0)

	go findFilesReadDir("1", filesDir, filesPatterns, fileNamesChannel)
	// go findFilesGlob("main", filesDir, filesPatterns, fileNamesChannel)

	//go findTest("1", "test1", fileNamesChannel)
	//go findTest("2", "test2", fileNamesChannel)

	go processReadFile("1", fileNamesChannel, eventsChannel)
	//go processTest("2", fileNamesChannel, eventsChannel)
	//go processTest("3", fileNamesChannel, eventsChannel)

	go outputWriteFile("1", eventsChannel)

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
