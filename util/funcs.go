package util

import (
	"log"
	"time"
)

func TimeMe(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s: %s", name, elapsed)
}
