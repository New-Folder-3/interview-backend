package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

var Log = log.New()

func ErrorPrinter(err error) {
	if err != nil {
		log.Errorln(err)
		stackTrace := fmt.Sprintf("%+v", err)
		lines := strings.Split(stackTrace, "\n")
		for _, line := range lines {
			if strings.Contains(line, "interview") {
				log.Errorln(line)
			}
		}
	}
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}

func StructPrinter(s interface{}, indent int) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// 指针处理
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		fmt.Printf("%s%v\n", spaces(indent), val.Interface())
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		fmt.Printf("%s%s: ", spaces(indent), field.Name)

		// 递归处理嵌套结构体
		switch value.Kind() {
		case reflect.Struct:
			fmt.Println()
			StructPrinter(value.Interface(), indent+2)
		case reflect.Ptr:
			if value.IsNil() {
				fmt.Println("<nil>")
			} else {
				fmt.Println()
				StructPrinter(value.Interface(), indent+2)
			}
		default:
			fmt.Printf("%v\n", value.Interface())
		}
	}
}
