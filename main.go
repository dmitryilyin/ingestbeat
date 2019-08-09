package main

import (
	"os"

	"github.com/dmitryilyin/ingestbeat/cmd"

	_ "github.com/dmitryilyin/ingestbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
