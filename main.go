package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func getCrontabs() ([]string) {
	var crontabs []string
	crontabs = append(crontabs, "test_data/crontab")
	return crontabs
}

func processCrontabFile(crontab string) {
	log.Printf("%s\n", crontab)
	crontab_entries, err := ioutil.ReadFile(crontab)
	if err != nil {
		log.Printf("WARNING: %s missing", crontab)
		return
	}

	for _, crontab_entry := range strings.Split(string(crontab_entries), "\n") {
		if len(crontab_entry) > 0 {
			processCrontabEntry(crontab_entry)
		}
	}
}

func logToStdout() {
	loggingFilename := "/dev/stdout"
	f, err := os.OpenFile(loggingFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}

func main() {
	//logToStdout()
	for _, crontab := range getCrontabs() {
		processCrontabFile(crontab)
	}
}