package main

import (
	"testing"
	"time"
)

func TestIsCurrentDayOfWeek_Monday(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Mon"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_Monday")
	}
}

func TestIsCurrentDayOfWeek_Wildcard(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "*"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_Wildcard")
	}
}

func TestIsCurrentDayOfWeek_ListWithWildcard(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Tue,Wed,*"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_ListWithWildcard")
	}
}

func TestIsCurrentDayOfWeek_ListWithRangeWithWildcard(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Tue,Wed,*-Sat"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_ListWithRangeWithWildcard")
	}
}

func TestIsCurrentDayOfWeek_Sunday(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "0"
	now :=time.Date(2019, 2,24,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_Sunday")
	}
}

func TestIsCurrentDayOfWeek_TuesdayInRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Sun-MON,TUE-Thu"
	now :=time.Date(2019, 2,26,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_TuesdayInRange")
	}
}

func TestIsCurrentDayOfWeek_WednesdayInRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Sun-MON,2-5"
	now :=time.Date(2019, 2,27,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_WednesdayInRange")
	}
}

func TestIsCurrentDayOfWeek_ThursdayNotInRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Sun-Wed"
	now :=time.Date(2019, 2,28,10,0,0,0,time.UTC)
	if crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_ThursdayNotInRange")
	}
}

func TestIsCurrentDayOfWeek_ThursdayInRange_Step(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Sun/2"
	now :=time.Date(2019, 2,28,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_ThursdayInRange_Step")
	}
}

func TestIsCurrentDayOfWeek_ThursdayNotInRange_Step(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Mon/2"
	now :=time.Date(2019, 2,28,10,0,0,0,time.UTC)
	if crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_ThursdayNotInRange_Step")
	}
}

func TestIsCurrentDayOfWeek_SaturdayInRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Sun-Sat"
	now :=time.Date(2019, 3,2,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_SaturdayInRange")
	}
}

func TestIsCurrentDayOfWeek_SaturdayInList(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfWeek = "Sun-Tue,3-4,6"
	now :=time.Date(2019, 3,2,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfWeek(now) {
		t.Errorf("Failed TestIsCurrentDayOfWeek_SaturdayInList")
	}
}

func TestIsCurrentMonth_February(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "Feb"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_February")
	}
}

func TestIsCurrentMonth_Wildcard(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "*"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_Wildcard")
	}
}

func TestIsCurrentMonth_StepWithWildcard(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "*/3"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_Wildcard")
	}
}

func TestIsCurrentMonth_NotMarch(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "Mar"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_NotMarch")
	}
}

func TestIsCurrentMonth_JuneInRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "May-Sep"
	now :=time.Date(2019, 6,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_JuneInRange")
	}
}

func TestIsCurrentMonth_JuneNotInRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "Jul-Dec"
	now :=time.Date(2019, 6,25,10,0,0,0,time.UTC)
	if crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_JuneNotInRange")
	}
}

func TestIsCurrentMonth_JanInList(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "3,5,7,Oct,1-Feb"
	now :=time.Date(2019, 1,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_JanInList")
	}
}

func TestIsCurrentMonth_DecInStep(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "Jan,9/3"
	now :=time.Date(2019, 12,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_DecInStep")
	}
}

func TestIsCurrentMonth_NovNotInStep(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.month = "2,Sep/3"
	now :=time.Date(2019, 11,25,10,0,0,0,time.UTC)
	if crontab_entry.isCurrentMonth(now) {
		t.Errorf("Failed TestIsCurrentMonth_NovNotInStep")
	}
}

func TestIsCurrentDayOfMonth_25(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfMonth = "25"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfMonth(now) {
		t.Errorf("Failed TestIsCurrentDayOfMonth_25")
	}
}

func TestIsCurrentDayOfMonth_InList(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfMonth = "1-5,11/2"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfMonth(now) {
		t.Errorf("Failed TestIsCurrentDayOfMonth_InList")
	}
}

func TestIsCurrentDayOfMonth_Wildcard(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfMonth = "*"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfMonth(now) {
		t.Errorf("Failed TestIsCurrentDayOfMonth_Wildcard")
	}
}

func TestIsCurrentDayOfMonth_WildcardInList(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfMonth = "1-*"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfMonth(now) {
		t.Errorf("Failed TestIsCurrentDayOfMonth_WildcardInList")
	}
}

func TestIsCurrentDayOfMonth_NotInRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfMonth = "1-24"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if crontab_entry.isCurrentDayOfMonth(now) {
		t.Errorf("Failed TestIsCurrentDayOfMonth_NotInRange")
	}
}

func TestIsCurrentDayOfMonth_Step(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfMonth = "2/2"
	now :=time.Date(2019, 3,30,10,0,0,0,time.UTC)
	if !crontab_entry.isCurrentDayOfMonth(now) {
		t.Errorf("Failed TestIsCurrentDayOfMonth_Step")
	}
}

func TestIsCurrentDayOfMonth_InvalidRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.dayOfMonth = "44-55"
	now :=time.Date(2019, 2,25,10,0,0,0,time.UTC)
	if crontab_entry.isCurrentDayOfMonth(now) {
		t.Errorf("Failed TestIsCurrentDayOfMonth_InvalidRange")
	}
}

func TestIsCurrentHour_10(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.hour = "10"
	now := time.Date(2019, 2, 25, 10, 0, 0, 0, time.UTC)
	if !crontab_entry.isCurrentHour(now) {
		t.Errorf("Failed TestIsCurrentHour_10")
	}
}

func TestIsCurrentHour_Not10(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.hour = "10"
	now := time.Date(2019, 2, 25, 11, 0, 0, 0, time.UTC)
	if crontab_entry.isCurrentHour(now) {
		t.Errorf("Failed TestIsCurrentHour_Not10")
	}
}

func TestIsCurrentHour_ListOfMatches(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.hour = "4,8,10,14,44"
	now := time.Date(2019, 2, 25, 10, 0, 0, 0, time.UTC)
	if !crontab_entry.isCurrentHour(now) {
		t.Errorf("Failed TestIsCurrentHour_ListOfMatches")
	}
}

func TestIsCurrentHour_InvalidRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.hour = "30-40"
	now := time.Date(2019, 2, 25, 10, 0, 0, 0, time.UTC)
	if crontab_entry.isCurrentHour(now) {
		t.Errorf("Failed TestIsCurrentHour_InvalidRange")
	}
}



func TestIsCurrentMinute_22(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.minute = "22"
	now := time.Date(2019, 2, 25, 10, 22, 0, 0, time.UTC)
	if !crontab_entry.isCurrentMinute(now) {
		t.Errorf("Failed TestIsCurrentMinute_22")
	}
}

func TestIsCurrentMinute_Not22(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.minute = "21"
	now := time.Date(2019, 2, 25, 11, 22, 0, 0, time.UTC)
	if crontab_entry.isCurrentMinute(now) {
		t.Errorf("Failed TestIsCurrentMinute_Not22")
	}
}

func TestIsCurrentMinute_ListOfMatches(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.minute = "4,8,10,14,44"
	now := time.Date(2019, 2, 25, 10, 14, 0, 0, time.UTC)
	if !crontab_entry.isCurrentMinute(now) {
		t.Errorf("Failed TestIsCurrentHour_ListOfMatches")
	}
}

func TestIsCurrentMinute_InvalidRange(t *testing.T) {
	crontab_entry := CrontabEntryFactory()
	crontab_entry.minute = "30-"
	now := time.Date(2019, 2, 25, 10, 0, 0, 0, time.UTC)
	if crontab_entry.isCurrentMinute(now) {
		t.Errorf("Failed TestIsCurrentMinute_InvalidRange")
	}
}