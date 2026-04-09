package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "invalid now", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dateStr == "" {
		http.Error(w, "date is required", http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat for \"d\" rule")
		}
		number, err := strconv.Atoi(parts[1])
		if err != nil || number < 1 || number > 400 {
			return "", errors.New("invalid number of days")
		}
		for {
			date = date.AddDate(0, 0, number)
			if date.After(now) {
				return date.Format("20060102"), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid repeat for \"y\" rule")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format("20060102"), nil
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat for \"w\" rule")
		}
		date, err = nextWeekday(date, now, parts[1])
		return date.Format("20060102"), err

	case "m":
		if len(parts) != 2 && len(parts) != 3 {
			return "", errors.New("invalid repeat for \"m\" rule")
		}
		monthsPart := ""
		if len(parts) == 3 {
			monthsPart = parts[2]
		}
		date, err = nextMonthday(date, now, parts[1], monthsPart)
		return date.Format("20060102"), err

	default:
		return "", errors.New("unsupported repeat rule")
	}
}

// проверяет, что дата больше now
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func nextWeekday(date time.Time, now time.Time, days string) (time.Time, error) {
	var weekdays [8]bool
	parts := strings.Split(days, ",")
	for _, num := range parts {
		day, err := strconv.Atoi(num)
		if err != nil || day < 1 || day > 7 {
			return time.Time{}, errors.New("invalid weekday")
		}
		weekdays[day] = true
	}

	for {
		dow := int(date.Weekday())
		if dow == 0 { // меняем номер воскресенья с 0 на 7
			dow = 7
		}
		if weekdays[dow] {
			if afterNow(date, now) {
				return date, nil
			}
		}
		date = date.AddDate(0, 0, 1)
	}
}

func nextMonthday(date time.Time, now time.Time, daysPart string, monthsPart string) (time.Time, error) {
	var dayArr [32]bool
	var monthArr [13]bool

	days := strings.Split(daysPart, ",")
	for _, num := range days {
		day, err := strconv.Atoi(num)
		if err != nil || day < -2 || day > 31 {
			return time.Time{}, errors.New("invalid day")
		}
		dayArr[day] = true
	}

	if len(monthsPart) != 0 {
		months := strings.Split(monthsPart, ",")
		for _, num := range months {
			month, err := strconv.Atoi(num)
			if err != nil || month < 1 || month > 12 {
				return time.Time{}, errors.New("invalid month")
			}
			monthArr[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			monthArr[i] = true
		}
	}

	for {
		// последний и предпоследний день месяца
		lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
		secondLastDay := lastDay - 1

		// попадает ли текущая дата из date в массив дней
		validDay := dayArr[date.Day()]
		for _, num := range days {
			day, _ := strconv.Atoi(num)
			if (day == -1 && date.Day() == lastDay) || (day == -2 && date.Day() == secondLastDay) {
				validDay = true
				break
			}
		}

		// попадает ли меясц из date в массив месцев
		validMonth := monthArr[date.Month()]

		if validDay && validMonth && date.After(now) {
			return date, nil
		}

		date = date.AddDate(0, 0, 1)
	}
}
