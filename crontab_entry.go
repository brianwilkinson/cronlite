package main

import (
	"log"
	"strings"
	"time"
)

const (
	minute_or_non_standard = iota
	minute
	hour
	dayofmonth
	month
	dayofweek
	command
	end
)

type CrontabEntry struct {
	minute      string
	hour        string
	dayOfMonth  string
	month       string
	dayOfWeek   string
	nonStandard string
	command     string
}

func CrontabEntryFactory() (entry *CrontabEntry) {
	return &CrontabEntry{
		minute:      "",
		hour:        "",
		dayOfMonth:  "",
		month:       "",
		dayOfWeek:   "",
		nonStandard: "",
		command:     "",
	}
}

func processCrontabEntry(crontab_entry_line string) () {
	crontab_entry := parseCrontabEntryLine(crontab_entry_line)
	log.Printf("TEST %v", crontab_entry)
	if crontab_entry.nonStandard != "" {
		runNonStandardCrontabEntry(crontab_entry)
	} else {
		runCrontabEntry(crontab_entry)
	}
}

func parseCrontabEntryLine(crontab_entry_line string) (*CrontabEntry) {
	log.Printf("Entry: %s\n", crontab_entry_line)
	crontab_entry := CrontabEntryFactory()
	state := minute_or_non_standard
	crontab_entry_parts := strings.Split(crontab_entry_line, " ")
	for index, crontab_entry_part := range crontab_entry_parts {
		switch state {
		case minute_or_non_standard:
			if strings.Index(crontab_entry_part, "@") == 0 {
				crontab_entry.nonStandard = crontab_entry_part
				state = command
			} else {
				crontab_entry.minute = crontab_entry_part
				state = hour
			}
		case hour:
			crontab_entry.hour = crontab_entry_part
			state = dayofmonth
		case dayofmonth:
			crontab_entry.dayOfMonth = crontab_entry_part
			state = month
		case month:
			crontab_entry.month = crontab_entry_part
			state = dayofweek
		case dayofweek:
			crontab_entry.dayOfWeek = crontab_entry_part
			state = command
		case command:
			crontab_entry.command = strings.Join(crontab_entry_parts[index:], " ")
			state = end
			break
		}
	}

	return crontab_entry
}

func (crontabEntry *CrontabEntry) isInDayOfWeek() (bool) {
	weekday := time.Now().Weekday().String()[0:3]
	log.Print(weekday)

	return true
}

func (crontabEntry *CrontabEntry) isCurrentMonth() (bool) {
	return true
}

func (crontabEntry *CrontabEntry) isCurrentDayOfMonth() (bool) {
	return true
}

func (crontabEntry *CrontabEntry) isCurrentHour() (bool) {
	return true
}

func (crontabEntry *CrontabEntry) isCurrentMinute() (bool) {
	return true
}

func runCrontabEntry(crontab_entry *CrontabEntry) {
	state := dayofweek
	for state != end {
		switch state {
		case dayofweek:
			if crontab_entry.dayOfWeek == "*" || crontab_entry.isInDayOfWeek() {
				state = month
			} else {
				log.Print("Day of week not matched")
				state = end
			}
		case month:
			if crontab_entry.month == "*" || crontab_entry.isCurrentMonth() {
				state = dayofmonth
			} else {
				log.Print("Month not matched")
				state = end
			}
		case dayofmonth:
			if crontab_entry.dayOfMonth == "*" || crontab_entry.isCurrentDayOfMonth() {
				state = hour
			} else {
				log.Print("Day of month not matched")
				state = end
			}
		case hour:
			if crontab_entry.hour == "*" || crontab_entry.isCurrentHour() {
				state = minute
			} else {
				log.Print("Hour not matched")
				state = end
			}
		case minute:
			if crontab_entry.minute == "*" || crontab_entry.isCurrentMinute() {
				log.Print("Command would run")
			} else {
				log.Print("Minute not matched")
			}
			state = end
		}
	}
}

func runNonStandardCrontabEntry(crontab_entry *CrontabEntry) {

}