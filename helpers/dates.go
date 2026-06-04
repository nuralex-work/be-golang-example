package helpers

import (
	"log"
	"strings"
	"time"
)

func GetFullDate(year int, month int, date int) (res []string) {
	log.Println("tahun: ", year)
	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	// Define the end date (one year from the start date)
	endDate := startDate.AddDate(1, 0, 0) // Add 1 year, 0 months, 0 days

	// Loop through each day from startDate to endDate
	for current := startDate; current.Before(endDate); current = current.AddDate(0, month, date) {
		if month > 0 {
			res = append(res, current.Format("01-2006"))
		} else {
			res = append(res, current.Format("2006-01-02"))
		}
	}
	return
}
func ParseToTime(dates string) (time.Time, error) {
	tgl, err := time.Parse(time.RFC3339, dates)
	return tgl, err
}
func ParseToTimeServer(dates string) (time.Time, error) {

	newdate := strings.Split(dates, "-")
	var arrdate []string
	for i := len(newdate); i > 0; i-- {
		arrdate = append(arrdate, newdate[i-1])
	}
	strdate := strings.Join(arrdate, "-")
	tgl, err := time.Parse("2006-01-02", strdate)
	return tgl, err
}
