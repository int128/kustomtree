package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/int128/sortmanifest/pkg/cmd"
)

func run(osArgs []string) error {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := f.Parse(osArgs[1:]); err != nil {
		return fmt.Errorf("wrong argument: %w", err)
	}
	if f.NArg() < 1 {
		return fmt.Errorf("you need to set at least 1 path")
	}
	for _, path := range f.Args() {
		if err := cmd.Run(path); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("error: %s", err)
	}
}
