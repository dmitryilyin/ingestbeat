package beater

import (
	"bufio"
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func (bt *Ingestbeat) processReadBuffer(instance string, fileNamesChannel <-chan string, b *beat.Beat) {
	if instance == "" {
		instance = "0"
	}
	logp.Debug("process", "ProcessReadFile[%s]: start", instance)

	for filePath := range fileNamesChannel {
		logp.Debug("process", "ProcessReadFile[%s]: receive \"%s\"", instance, filePath)

		fileHandler, err := os.Open(filePath)
		if err != nil {
			logp.Warn("ProcessReadFile[%s]: could not open: \"%s\", skipping", instance, filePath)
		}
		defer file.Close()

		fileScanner := bufio.NewScanner(fileHandler)

		for fileScanner.Scan() {
			fmt.Println(fileScanner.Text())
		}

		logp.Debug("process", "ProcessReadFile[%s]: sending %d bytes", instance, len(content))
		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"message": string(content),
			},
		}
		bt.client.Publish(event)
		logp.Debug("process", "ProcessReadFile[%s]: marking \"%s\"", instance, filePath)
		bt.processMarkDelete(filePath)
	}
}

func (bt *Ingestbeat) processReadFile(instance string, fileNamesChannel <-chan string, b *beat.Beat) {
	if instance == "" {
		instance = "0"
	}
	logp.Debug("process", "ProcessReadFile[%s]: start", instance)

	for filePath := range fileNamesChannel {
		logp.Debug("process", "ProcessReadFile[%s]: receive \"%s\"", instance, filePath)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			logp.Error(err)
		}
		logp.Debug("process", "ProcessReadFile[%s]: sending %d bytes", instance, len(content))
		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"message": string(content),
			},
		}
		bt.client.Publish(event)
		logp.Debug("process", "ProcessReadFile[%s]: marking \"%s\"", instance, filePath)
		bt.processMarkDelete(filePath)
	}
}

func (bt *Ingestbeat) processMarkRename(filePath string) {
	newFilePath := filePath + ".processed"
	logp.Debug("process", "ProcessMarkRename: \"%s\" -> \"%s\"", filePath, newFilePath)
	err := os.Rename(filePath, newFilePath)
	if err != nil {
		logp.Error(err)
	}
}

func (bt *Ingestbeat) processMarkDelete(filePath string) {
	logp.Debug("process", "ProcessMarkDelete: \"%s\"", filePath)
	err := os.Remove(filePath)
	if err != nil {
		logp.Error(err)
	}
}
