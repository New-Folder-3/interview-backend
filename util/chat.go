package util

import (
	"strconv"
	"strings"
)

func ChatStringToFloatSlice(str string) []float64 {
	if str == "" {
		return nil
	}
	var result []float64
	items := strings.Split(str, " ")
	for _, item := range items {
		f, err := strconv.ParseFloat(item, 64)
		if err == nil {
			result = append(result, f)
		}
	}
	return result
}

func ChatStringToStringSlice(str string) []string {
	if str == "" {
		return nil
	}
	var result []string
	items := strings.Split(str, " ")
	for _, item := range items {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}
