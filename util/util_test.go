package util

import (
	"fmt"
	"testing"
)

func TestStringListToDB(t *testing.T) {
	test := []string{"a", "b", "c"}
	fmt.Println(StringListToDB(test))
}
