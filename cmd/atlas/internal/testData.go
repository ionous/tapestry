package internal

import (
	"database/sql"
	_ "embed"
)

func CreateTestData(db *sql.DB) (err error) {
	_, err = db.Exec(testData)
	return
}

//go:embed testData.sql
var testData string
