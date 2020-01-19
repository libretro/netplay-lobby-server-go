package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GetPostgreDB returns a postgresql interface using GORM.
func GetPostgreDB(connection string) (*gorm.DB, error) {
	return gorm.Open("postgres", connection)
}
