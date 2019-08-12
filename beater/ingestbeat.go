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

	var fileNamesChannel = make(chan string, 0)

	filesDir := "lib/test"
	filesPatterns := []string{"*.log"}

	logp.Info("Starting finders")
	go bt.findFilesReadDir("0", filesDir, filesPatterns, fileNamesChannel)

	logp.Info("Starting processors")
	go bt.processReadFile("0", fileNamesChannel, b)

	for {
		logp.Info("Starting main loop")
		select {
		case <-bt.done:
			return nil
		}
	}
}

// Stop stops ingestbeat.
func (bt *Ingestbeat) Stop() {
	_ = bt.client.Close()
	close(bt.done)
}
