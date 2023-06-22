package util

import "strings"

func BooleanToYesNo(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func YesNoToBoolean(value string) bool {
	return strings.ToLower(value) == "yes"
}
