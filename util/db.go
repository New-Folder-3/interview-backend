package util

import (
	"encoding/json"
	"fmt"
	"interview/internal/conf"
	"strings"
)

func DbGetColumnName(name string) string {
	if conf.Conf.Database.Type == "postgres" {
		return fmt.Sprintf(`"%s"`, name)
	}
	return fmt.Sprintf("`%s`", name)
}

func StringListToDB(list []string) string {
	if len(list) == 0 {
		return ""
	}
	result := ""
	for i, item := range list {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%s", item)
	}
	return result
}

func DBToStringList(dbList string) []string {
	if dbList == "" {
		return nil
	}
	var result []string
	items := strings.Split(dbList, ",")
	for _, item := range items {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func RemoveFromDBList(dbList string, item string) string {
	if dbList == "" {
		return ""
	}
	items := strings.Split(dbList, ",")
	var result []string
	for _, i := range items {
		if i != item {
			result = append(result, i)
		}
	}
	return StringListToDB(result)
}

func MapToStruct(m map[string]interface{}, out interface{}) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}
