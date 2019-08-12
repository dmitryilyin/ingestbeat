package beater

import (
	"github.com/elastic/beats/libbeat/logp"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func (bt *Ingestbeat) findCheckSkip(fileInfo os.FileInfo) bool {
	logp.Debug("find", "FindCheckMarked: " + fileInfo.Name())
	if fileInfo.IsDir() {
		return false
	}
	if fileInfo.Size() == 0 {
		return false
	}
	return true
}

func (bt *Ingestbeat) findCheckMatch(fileInfo os.FileInfo, filePattern string) bool {
	logp.Debug("find", "FindCheckMatch: " + fileInfo.Name())
	if fileInfo.IsDir() {
		return false
	}
	matched, err := filepath.Match(filePattern, fileInfo.Name())
	if err != nil {
		logp.Error(err)
		return false
	}
	return matched
}

func (bt *Ingestbeat) findFilesReadDir(instance string, fileDirectory string, filePatterns []string, fileNamesChannel chan<- string) {
	if instance == "" {
		instance = "0"
	}
	logp.Debug("find", "FindFilesReadDir[%s]: start", instance)
	for {
		logp.Debug("find", "FindFilesReadDir[%s]: scan \"%s\"", instance, fileDirectory)
		fileInfoList, err := ioutil.ReadDir(fileDirectory)
		if err != nil {
			logp.Error(err)
		}
		for _, fileInfo := range fileInfoList {
			for _, filePattern := range filePatterns {
				if !bt.findCheckSkip(fileInfo) {
					continue
				}
				if !bt.findCheckMatch(fileInfo, filePattern) {
					continue
				}
				fileFullPath := fileDirectory + string(os.PathSeparator) + fileInfo.Name()
				logp.Debug("find", "FindFilesReadDir[%s]: send \"%s\"", instance, fileFullPath)
				fileNamesChannel <- fileFullPath
			}
		}

		time.Sleep(3 * time.Second)
	}
}
