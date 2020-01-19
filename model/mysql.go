package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// GetMysqlDB returns a mysql interface using GORM.
func GetMysqlDB(connection string) (*gorm.DB, error) {
	return gorm.Open("mysql", connection)
}
