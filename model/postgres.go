package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GetPostgreDB returns a postgresql interface using GORM.
func GetPostgreDB(host string, port uint, username string, password string, dbname string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, username, password, dbname)
	return gorm.Open("postgres", connectionString)
}
