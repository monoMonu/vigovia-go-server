package utils

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// FormatDate converts a string date to the desired format.
func FormatDate(date string) string {
	// Parse the date string into a time.Time object
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return ""
	}
	return parsedDate.Format("02 Jan, 2006")
}

// CalculateNights calculates the number of nights between two date strings.
func CalculateNights(departureDate string, returnDate string) int {
	// Parse both dates into time.Time objects
	departure, err1 := time.Parse("2006-01-02", departureDate)
	returnDateParsed, err2 := time.Parse("2006-01-02", returnDate)
	if err1 != nil || err2 != nil {
		log.Println("Error parsing dates:", err1, err2)
		return 0
	}
	return int(returnDateParsed.Sub(departure).Hours() / 24)
}

// ConvertStringToTime converts a string date into a time.Time object.
func ConvertStringToTime(date string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	return parsedDate, nil
}

// SanitizeFileName cleans up a filename by removing invalid characters.
func SanitizeFileName(name string) string {
	res := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			res += string(r)
		} else if r == ' ' || r == '_' || r == '-' {
			res += "_"
		}
	}
	return res
}

func FormatDuration(duration int, timeSlot string) string {
	if duration > 0 {
		if duration >= 60 {
			hours := duration / 60
			minutes := duration % 60
			if minutes == 0 {
				if hours == 1 {
					return "1 Hour"
				}
				return fmt.Sprintf("%d Hours", hours)
			}
			return fmt.Sprintf("%d:%02d Hours", hours, minutes)
		}
		return fmt.Sprintf("%d Minutes", duration)
	}

	// If no duration specified, try to infer from time slot
	timeSlot = strings.ToLower(timeSlot)
	if strings.Contains(timeSlot, "morning") || strings.Contains(timeSlot, "afternoon") || strings.Contains(timeSlot, "evening") {
		return "2-3 Hours"
	}
	if strings.Contains(timeSlot, "full day") || strings.Contains(timeSlot, "all day") {
		return "Full Day"
	}
	if strings.Contains(timeSlot, "half day") {
		return "Half Day"
	}

	// Default
	return "2-3 Hours"
}
