package utils

import (
	"github.com/sirupsen/logrus"
	"time"
)

func TimeTrack(start time.Time, name string, log logrus.StdLogger) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
