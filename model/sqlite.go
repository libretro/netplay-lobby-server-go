package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// GetSqliteDB returns a sqlite3 interface using GORM.
// A filepath of ":memory:" provides a in memory database.
func GetSqliteDB(filepath string) (*gorm.DB, error) {
	return gorm.Open("sqlite3", filepath)
}
