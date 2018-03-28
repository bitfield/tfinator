package main

import (
	"fmt"
	"os"

	"github.com/carezone/tfinator"
)

func main() {
	for _, path := range os.Args[1:] {
		fmt.Println(tfinator.PlanStats(path))
	}
}
