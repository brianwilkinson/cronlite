package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
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
	if crontab_entry.nonStandard != "" {
		go runNonStandardCrontabEntry(crontab_entry)
	} else {
		go runCrontabEntry(crontab_entry)
	}
}

func parseCrontabEntryLine(crontab_entry_line string) (*CrontabEntry) {
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

func (crontabEntry *CrontabEntry) isCurrentDayOfWeek(now time.Time) (bool) {
	// maps are not ordered, so we use a list here...
	daylist := [7]string {"SUN","MON","TUE","WED","THU","FRI","SAT"}

	// and then work out the current_weekday_num based on the list index
	current_weekday_txt := strings.ToUpper(now.Weekday().String()[0:3])
	current_weekday_num := "0"
	for index, val := range daylist {
		if val == current_weekday_txt {
			current_weekday_num = strconv.Itoa(index)
			break
		}
	}

	for _, dayranges := range strings.Split(strings.ToUpper(crontabEntry.dayOfWeek), ",") {
		dayrange := strings.Split(dayranges, "-")
		if len(dayrange) > 1 {
			start_day := dayrange[0]
			end_day := dayrange[1]
			if start_day == "*" || end_day == "*" {
				return true
			}

			valid := false
			for index, daytxt := range daylist {
				daynum := strconv.Itoa(index)
				if !valid {
					if start_day == daytxt || start_day == daynum {
						valid = true
					}
				}
				if valid {
					if daytxt == current_weekday_txt || daynum == current_weekday_num {
						return true
					}
				}
				if end_day == daytxt || end_day == daynum {
					break
				}
			}
		} else {
			daystep := strings.Split(dayrange[0], "/")
			if len(daystep) == 2 {
				// Handle steps e.g. Tue/3 == every 3 days starting tuesday == Tue,Fri
				day_val := daystep[0]
				if day_val == "*" {
					return true
				}
				day_vals := []string{}

				step_val, err := strconv.Atoi(daystep[1])
				if err == nil {
					for index, daytxt := range daylist {
						daynum := strconv.Itoa(index)
						if day_val == daytxt || day_val == daynum {
							day_vals = append(day_vals, daytxt)
							for {
								index += step_val
								if index < len(daylist) {
									day_vals = append(day_vals, daylist[index])
								} else {
									break
								}
							}
							for _, valid_day_val := range day_vals {
								if valid_day_val == current_weekday_txt {
									return true
								}
							}
							break
						}
					}
				}
			} else if len(daystep) == 1 {
				// Just a single day of week e.g. Fri or 5
				ct_day := dayrange[0]
				if ct_day == "*" || ct_day == current_weekday_txt || ct_day == current_weekday_num {
					return true
				}
			}
		}
	}

	return false
}

func (crontabEntry *CrontabEntry) isCurrentMonth(now time.Time) (bool) {
	// maps are not ordered, so we use a list here...
	monthlist := [12]string {"JAN","FEB","MAR","APR","MAY","JUN","JUL","AUG","SEP","OCT","NOV","DEC"}

	// and then work out the current_month_num based on the list index
	current_month_txt := strings.ToUpper(now.Month().String()[0:3])
	current_month_num := "1"
	for index, val := range monthlist {
		if val == current_month_txt {
			current_month_num = strconv.Itoa(index + 1)
			break
		}
	}

	for _, monthranges := range strings.Split(strings.ToUpper(crontabEntry.month), ",") {
		monthrange := strings.Split(monthranges, "-")
		if len(monthrange) > 1 {
			// We have a range of months e.g. Mar-Jun == Mar,Apr,May,Jun
			start_month := monthrange[0]
			end_month := monthrange[1]
			if start_month == "*" || end_month == "*" {
				return true
			}
			valid := false
			for index, monthtxt := range monthlist {
				monthnum := strconv.Itoa(index + 1)
				if !valid {
					if start_month == monthtxt || start_month == monthnum {
						valid = true
					}
				}
				if valid {
					if monthtxt == current_month_txt || monthnum == current_month_num {
						return true
					}
				}
				if end_month == monthtxt || end_month == monthnum {
					break
				}
			}
		} else {
			monthstep := strings.Split(monthrange[0], "/")
			if len(monthstep) == 2 {
				// Handle steps e.g. Feb/3 == every 3 months starting February == Feb,May,Aug,Nov
				month_val := monthstep[0]
				if month_val == "*" {
					return true
				}
				month_vals := []string{}

				step_val, err := strconv.Atoi(monthstep[1])
				if err == nil {
					for index, monthtxt := range monthlist {
						monthnum := strconv.Itoa(index + 1)
						if month_val == monthtxt || month_val == monthnum {
							month_vals = append(month_vals, monthtxt)
							for {
								index += step_val
								if index < len(monthlist) {
									month_vals = append(month_vals, monthlist[index])
								} else {
									break
								}
							}
							for _, valid_month_val := range month_vals {
								if valid_month_val == current_month_txt {
									return true
								}
							}
							break
						}
					}
				}
			} else if len(monthstep) == 1 {
				// Just a single month e.g. Mar
				ct_month := monthrange[0]
				if ct_month == "*" || ct_month == current_month_txt || ct_month == current_month_num {
					return true
				}
			}
		}
	}

	return false
}

func getIntVal(val string, minimum int, maximum int) (int, error) {
	int_val, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	if int_val < minimum {
		return 0, errors.New("Below minimum value")
	}
	if int_val > maximum {
		return 0, errors.New("Above maximum value")
	}
	return int_val, err
}

func isNumericalRangeMatch(ct_val string, minimum, maximum, dt_val int) (bool) {
	for _, ct_ranges := range strings.Split(strings.ToUpper(ct_val), ",") {
		ct_range := strings.Split(ct_ranges, "-")
		if len(ct_range) > 1 {
			// We have a range of days e.g. 5-20 == 5th to 20th
			start_txt := ct_range[0]
			end_txt := ct_range[1]
			if start_txt == "*" || end_txt == "*" {
				return true
			}

			start_val, err := getIntVal(start_txt, minimum, maximum)
			if err != nil {
				return false
			}
			end_val, err := getIntVal(end_txt, minimum, maximum)
			if err != nil {
				return false
			}

			if start_val <= dt_val && dt_val <= end_val {
				return true
			}
		} else {
			ct_step := strings.Split(ct_range[0], "/")
			if len(ct_step) == 2 {
				// Handle steps e.g. 23/2 == 23rd/25th/27th/29th/31st
				if ct_step[0] == "*" {
					return true
				}

				start_val, err := getIntVal(ct_step[0], minimum, maximum)
				if err != nil {
					return false
				}
				step_val, err := getIntVal(ct_step[1], minimum, maximum)
				if err != nil {
					return false
				}

				for val := start_val; val <= maximum; val += step_val {
					if val == dt_val {
						return true
					}
				}
			} else if len(ct_step) == 1 {
				// Just a single dom e.g. 11 (== 11th)
				ct_dom := ct_range[0]
				if ct_dom == "*" || ct_dom == strconv.Itoa(dt_val) {
					return true
				}
			}
		}
	}

	return false
}

func (crontabEntry *CrontabEntry) isCurrentDayOfMonth(now time.Time) (bool) {
	return isNumericalRangeMatch(crontabEntry.dayOfMonth, 1, 31, now.Day())
}

func (crontabEntry *CrontabEntry) isCurrentHour(now time.Time) (bool) {
	return isNumericalRangeMatch(crontabEntry.hour, 0, 23, now.Hour())
}

func (crontabEntry *CrontabEntry) isCurrentMinute(now time.Time) (bool) {
	return isNumericalRangeMatch(crontabEntry.minute, 0, 59, now.Minute())
}

func runCrontabEntry(crontab_entry *CrontabEntry) {
	state := dayofweek
	for state != end {
		switch state {
		case dayofweek:
			if crontab_entry.isCurrentDayOfWeek(time.Now()) {
				state = month
			} else {
				log.Print("Day of week not matched")
				state = end
			}
		case month:
			if crontab_entry.isCurrentMonth(time.Now()) {
				state = dayofmonth
			} else {
				log.Print("Month not matched")
				state = end
			}
		case dayofmonth:
			if crontab_entry.isCurrentDayOfMonth(time.Now()) {
				state = hour
			} else {
				log.Print("Day of month not matched")
				state = end
			}
		case hour:
			if crontab_entry.isCurrentHour(time.Now()) {
				state = minute
			} else {
				log.Print("Hour not matched")
				state = end
			}
		case minute:
			if crontab_entry.isCurrentMinute(time.Now()) {
				crontab_entry.runCommand()
			} else {
				log.Print("Minute not matched")
			}
			state = end
		}
	}
}

func (crontabEntry *CrontabEntry) runCommand() {
	var cmd *exec.Cmd

	args := []string{"-c", crontabEntry.command}
	cmd = exec.Command("bash", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	runerr := cmd.Start()
	if runerr != nil {
		log.Printf("Failed to run: %s", crontabEntry.command)
	}

	log.Printf("CMD [%s]", crontabEntry.command)
	exitVal := cmd.Wait()
	if exitVal != nil {
		log.Printf("Error running [%s] %v", crontabEntry.command, exitVal)
	} else {
		// log.Printf("OK: [%s]", crontabEntry.command)
	}

}

func runNonStandardCrontabEntry(crontab_entry *CrontabEntry) {
	log.Fatal("I do not yet support @yearly, @annually, @monthly,@weekly, @daily, @hourly, @reboot")
}