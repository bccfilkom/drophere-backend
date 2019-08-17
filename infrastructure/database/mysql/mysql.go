package mysql

import (
	"time"

	"github.com/jinzhu/gorm"
)

// New creates new GORM database instance
func New(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db = db.
		Set("gorm:table_options", "DEFAULT CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci ENGINE=InnoDB").
		Set("gorm:auto_preload", false)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Minute)

	// db = db.Set("gorm:auto_preload", true)

	return db, nil
}
