package util

import (
	"fmt"
	"strings"
	"time"
)

// GetTodayLastHour ...
func GetTodayLastHour() (lastTime string, err error) {

	present := time.Now()

	yr, mnth, day := present.Date()
	intMonth := int(mnth)

	// Getting the date output in string
	urDate := fmt.Sprintf(`%v-%v-%v`, yr, intMonth, day)

	urTime := fmt.Sprintf(`%d:%d:%d`, 23, 59, 59)

	lastTime = fmt.Sprintf(`%s %s`, urDate, urTime)

	return lastTime, nil

}

// GetTodayDate ...
func GetTodayDate() (todayDate string, err error) {
	present := time.Now()

	yr, mnth, day := present.Date()
	intMonth := int(mnth)

	// Getting the date output in string
	todayDate = fmt.Sprintf(`%v-%v-%v`, yr, intMonth, day)

	return todayDate, nil
}

// GetFilteredPresentTime ...
func GetFilteredPresentTime() (final string, err error) {

	b := strings.SplitAfterN(fmt.Sprintf("%v", time.Now()), ".", 2)
	//fmt.Printf("%q\n", b[0])
	final = strings.TrimSpace(strings.Replace(b[0], ".", "", -1))
	//fmt.Printf("%q\n", final)

	return

}

// ConvertTimeToString ...
func ConvertTimeToString(timeGiven time.Time) (final string, err error) {

	b := strings.SplitAfterN(fmt.Sprintf("%v", timeGiven), "+", 2)
	final = strings.TrimSpace(strings.Replace(b[0], "+", "", -1))
	final = strings.TrimSpace(strings.Replace(final, "0000-01-01", "", -1))

	return
}
