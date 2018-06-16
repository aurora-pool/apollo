package helpers

import (
	"os"
)

func GetDBurl() string {
	if dbUrl := os.Getenv("AURORA_URL"); len(dbUrl) > 1 {
		return dbUrl
	}

	return ""
}
