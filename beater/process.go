package beater

import (
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/beat"
)

func (bt *Ingestbeat) processReadFile(instance string, fileNamesChannel <-chan string) {
	if instance == "" {
		instance = "0"
	}
	logp.Info("ProcessReadFile[%s]: start", instance)

	for filePath := range fileNamesChannel {
		logp.Info("ProcessReadFile[%s]: receive \"%s\"", instance, filePath)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			logp.Panic(err)
		}
		logp.Info("ProcessReadFile[%s]: sending %d bytes", instance, len(content))

		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"message": content,
			},
		}
		bt.client.Publish(event)

		bt.processMarkDelete(filePath)
	}
}

func (bt *Ingestbeat) processMarkRename(filePath string) {
	newFilePath := filePath + ".processed"
	logp.Info("ProcessMarkRename: \"%s\" -> \"%s\"", filePath, newFilePath)
	err := os.Rename(filePath, newFilePath)
	if err != nil {
		logp.Panic(err)
	}
}

func (bt *Ingestbeat) processMarkDelete(filePath string) {
	logp.Info("ProcessMarkDelete: \"%s\"", filePath)
	err := os.Remove(filePath)
	if err != nil {
		logp.Panic(err)
	}
}
