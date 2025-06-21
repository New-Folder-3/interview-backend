package util

import (
	"fmt"
	"interview-backend/internal/conf"
)

func DbGetColumnName(name string) string {
	if conf.Conf.Database.Type == "postgres" {
		return fmt.Sprintf(`"%s"`, name)
	}
	return fmt.Sprintf("`%s`", name)
}
