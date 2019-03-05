package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func getCrontabs() ([]string) {
	var crontabs []string

	override_crontabs := os.Getenv("CRONTABS")
	if override_crontabs == "" {
		crontabs = append(crontabs, "/etc/crontab")
		globs := []string{"/var/spool/cron/crontabs/*", "/etc/cron.d/*"}
		for _, glob := range globs {
			files, err := filepath.Glob(glob)
			if err == nil {
				for _, file := range files {
					// log.Printf("CRONTAB: %v", file)
					crontabs = append(crontabs, file)
				}
			}
		}
	} else {
		for _, crontab := range strings.Split(override_crontabs, ":") {
			crontabs = append(crontabs, crontab)
		}
	}

	return crontabs
}

func processCrontabFile(crontab string) {
	// log.Printf("%s\n", crontab)
	crontab_entries, err := ioutil.ReadFile(crontab)
	if err != nil {
		log.Printf("WARNING: %s missing", crontab)
		return
	}

	re_envar := regexp.MustCompile(`^[A-Za-z]+[A-Za-z0-9]*=`)

	for _, crontab_entry := range strings.Split(string(crontab_entries), "\n") {
		if len(crontab_entry) > 0 {
			crontab_entry = strings.TrimSpace(crontab_entry)
			if strings.Index(crontab_entry, "#") == 0 {
				// log.Printf("Ignoring comment: %s\n", crontab_entry)
				continue
			} else if re_envar.MatchString(crontab_entry) {
				// log.Printf("I think this is an envar: %s", crontab_entry)
				envar := strings.Split(crontab_entry, "=")
				key := envar[0]
				val := strings.Trim(strings.Join(envar[1:], "="), `"`)
				err := os.Setenv(key, val)
				if err != nil {
					log.Printf("WARNING: Could not set envar %s\n", crontab_entry)
				}
			} else {
				processCrontabEntry(crontab_entry)
			}
		}
	}
}

func main() {
	for {
		secondsNow := time.Now().Second()
		secondsWait := 60 - secondsNow
		timer := time.NewTimer(time.Duration(secondsWait) * time.Second)

		select {
		case <-timer.C:
			for _, crontab := range getCrontabs() {
				go processCrontabFile(crontab)
			}
		}
	}
}