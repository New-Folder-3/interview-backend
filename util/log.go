package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

var Log = log.New()

func ErrorPrinter(err error) {
	if err != nil {
		log.Errorln(err)
		stackTrace := fmt.Sprintf("%+v", err)
		lines := strings.Split(stackTrace, "\n")
		for _, line := range lines {
			if strings.Contains(line, "interview") {
				log.Errorln(line)
			}
		}
	}
}
