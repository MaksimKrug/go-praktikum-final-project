package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	return date.Format(dateFormat) > now.Format(dateFormat)
}

func nextDays(date, now time.Time, args []string) (time.Time, error) {
	if len(args) != 1 {
		return time.Time{}, fmt.Errorf("days repeat not one arguments: %v", args)
	}

	nextDateParsed, err := strconv.Atoi(args[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("repeat date %v not parsed: %w", args, err)
	}
	if nextDateParsed < 1 || nextDateParsed > 400 {
		return time.Time{}, fmt.Errorf("repeat not in [1,400]: %d", nextDateParsed)
	}

	for {
		date = date.AddDate(0, 0, nextDateParsed)
		if afterNow(date, now) {
			break
		}
	}

	return date, nil
}

func nextYear(date, now time.Time, args []string) (time.Time, error) {
	// не уверен что нужна эта проверка и что лучше сделать её здесь
	if len(args) != 0 {
		return time.Time{}, fmt.Errorf("year repeat must have no arguments: %v", args)
	}

	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date, nil
}

func nextWeek(date, now time.Time, args []string) (time.Time, error) {
	if len(args) != 1 {
		return time.Time{}, fmt.Errorf("weeks repeat not one arguments: %v", args)
	}

	parts := strings.Split(args[0], ",")
	var days [7]bool
	for _, val := range parts {
		valParsed, err := strconv.Atoi(val)
		if err != nil {
			return time.Time{}, fmt.Errorf("weekday not parsed: %v", args)
		}
		if valParsed < 1 || valParsed > 7 {
			return time.Time{}, fmt.Errorf("weekday not in [1,7] range: %v", args)
		}

		// обработаем воскресенье как в time.Weekday
		if valParsed == 7 {
			days[0] = true
		} else {
			days[valParsed] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) && days[int(date.Weekday())] {
			break
		}
	}

	return date, nil
}

func hasPossibleMonthDate(
	days [32]bool,
	months [13]bool,
	lastDay, penultimateDay bool,
) bool {
	for month, allowed := range months {
		if month == 0 || !allowed {
			continue
		}

		if lastDay || penultimateDay {
			return true
		}

		maxDay := time.Date(
			2024, time.Month(month)+1, 0,
			0, 0, 0, 0, time.UTC,
		).Day()

		for day, allowed := range days {
			if allowed && day >= 1 && day <= maxDay {
				return true
			}
		}
	}

	return false
}

func nextMonth(date, now time.Time, args []string) (time.Time, error) {
	if len(args) < 1 || len(args) > 2 {
		return time.Time{}, fmt.Errorf("expected days and optional months")
	}

	var days [32]bool
	var months [13]bool
	var lastDay, penultimateDay bool

	// parse days
	parts := strings.Split(args[0], ",")
	for _, val := range parts {
		valParsed, err := strconv.Atoi(val)
		if err != nil {
			return time.Time{}, fmt.Errorf("day not parsed: %v", args)
		}
		switch {
		case valParsed == -1:
			lastDay = true
		case valParsed == -2:
			penultimateDay = true
		case valParsed >= 1 && valParsed <= 31:
			days[valParsed] = true
		default:
			return time.Time{}, fmt.Errorf("invalid day: %d", valParsed)
		}

	}

	// parse months
	if len(args) == 2 {
		parts = strings.Split(args[1], ",")
		for _, val := range parts {
			valParsed, err := strconv.Atoi(val)
			if err != nil {
				return time.Time{}, fmt.Errorf("month not parsed: %v", args)
			}
			if valParsed < 1 || valParsed > 12 {
				return time.Time{}, fmt.Errorf("month not in [1,12] range: %v", args)
			}
			months[valParsed] = true
		}
	} else {
		for i := range months {
			months[i] = true
		}
	}

	// check day + month
	if !hasPossibleMonthDate(days, months, lastDay, penultimateDay) {
		return time.Time{}, fmt.Errorf("no possible date for specified days and months")
	}

	// find next date
	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}
		if !months[int(date.Month())] {
			continue
		}

		last := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
		penultimate := last - 1
		matchesLast := lastDay && date.Day() == last
		matchesPenultimate := penultimateDay && date.Day() == penultimate
		if days[date.Day()] || matchesLast || matchesPenultimate {
			break
		}

	}

	return date, nil

}

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("dstart %s not parsed: %w", dstart, err)
	}
	repeatParsed := strings.Fields(repeat)
	if len(repeatParsed) == 0 {
		return "", fmt.Errorf("repeate is empty")
	}

	var newDate time.Time
	switch repeatParsed[0] {
	case "d":
		newDate, err = nextDays(date, now, repeatParsed[1:])
	case "y":
		newDate, err = nextYear(date, now, repeatParsed[1:])
	case "w":
		newDate, err = nextWeek(date, now, repeatParsed[1:])
	case "m":
		newDate, err = nextMonth(date, now, repeatParsed[1:])
	default:
		return "", fmt.Errorf("repeat type is not supported: %v", repeatParsed)
	}
	if err != nil {
		return "", err
	}

	return newDate.Format(dateFormat), err

}
