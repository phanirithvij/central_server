package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB the gorm database exposed globally to all modules
var DB *gorm.DB

// InitDB initializes the gorm database
func InitDB(debug bool) {
	// TODO postgres or sqlite selection from f8
	var err error
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if debug {
		DB = DB.Debug()
	}
}

// GetDB gets the gorm database
func GetDB() *gorm.DB {
	if DB == nil {
		InitDB(false)
	}
	return DB
}
