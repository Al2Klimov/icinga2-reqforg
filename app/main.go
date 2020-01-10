package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	entries, errRD := ioutil.ReadDir(".")
	if errRD != nil {
		log.WithFields(log.Fields{"dir": ".", "error": errRD.Error()}).Fatal("Couldn't list dir")
		return
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".plugin") {
			go loadPlugin(entry.Name())
		}
	}

	select {}
}
