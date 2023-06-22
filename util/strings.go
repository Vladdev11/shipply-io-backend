package util

import "strings"

func ToSnakeCase(str string) string {
	var result string
	var words []string
	var lastPos int
	str = strings.TrimSpace(str)
	str = strings.Replace(str, "-", "_", -1)
	str = strings.Replace(str, ".", "_", -1)
	str = strings.Replace(str, " ", "_", -1)

	for i, v := range str {
		if v >= 'A' && v <= 'Z' {
			if i > 0 {
				words = append(words, str[lastPos:i])
			}
			lastPos = i
		}
	}

	if len(words) == 0 {
		return strings.ToLower(str)
	}

	words = append(words, str[lastPos:])
	for k, word := range words {
		if k > 0 {
			result += "_"
		}
		result += strings.ToLower(word)
	}
	return result
}

func ContainsUpper(str string) bool {
	for _, v := range str {
		if v >= 'A' && v <= 'Z' {
			return true
		}
	}
	return false
}

func ContainsLower(str string) bool {
	for _, v := range str {
		if v >= 'a' && v <= 'z' {
			return true
		}
	}
	return false
}

func SafeDereferenceString(strPtr *string) string {
	if strPtr != nil {
		return *strPtr
	}
	return ""
}
