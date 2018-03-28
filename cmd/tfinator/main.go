package main

import (
	"fmt"
	"log"
	"os"

	"github.com/carezone/tfinator"
)

func main() {
	for _, path := range os.Args[1:] {
		s, err := tfinator.PlanStats(path)
		if err != nil {
			log.Fatalf("couldn't get plan stats on %s: %v", path, err)
		}
		fmt.Printf("%s: %+v\n", path, s)
	}
}
