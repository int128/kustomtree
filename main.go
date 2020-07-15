package main

import (
	"log"
	"os"

	"github.com/int128/kustomtree/pkg/cmd"
)

func init() {
	log.SetFlags(0)
}

func main() {
	if err := cmd.Run(os.Args); err != nil {
		log.Fatalf("error: %s", err)
	}
}
