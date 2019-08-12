package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/dmitryilyin/ingestbeat/config"
)

// Ingestbeat configuration.
type Ingestbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of ingestbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Ingestbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts ingestbeat.
func (bt *Ingestbeat) Run(b *beat.Beat) error {
	logp.Info("ingestbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	var fileNamesChannel chan string = make(chan string, 0)

	filesDir := "test"
	filesPatterns := []string{"*.log"}

	go bt.findFilesReadDir("0", filesDir, filesPatterns, fileNamesChannel)

	go bt.processReadFile("0", fileNamesChannel)

	for {
		select {
		case <-bt.done:
			return nil
		}


	}
}

// Stop stops ingestbeat.
func (bt *Ingestbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
