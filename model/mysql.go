package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// GetMysqlDB returns a mysql interface using GORM.
func GetMysqlDB(host string, username string, password string, dbname string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, dbname)
	return gorm.Open("mysql", connectionString)
}
