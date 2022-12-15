package main

import (
  "fmt"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func gormConnect() (*gorm.DB, error) {
  dbHost := "127.0.0.1"
  dbPort := "3306"
  dbName := "go_scraper_dev"
  dbUser := "user"
  dbPassword := "password"

  dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if err != nil {
    return nil, fmt.Errorf("db connection error : %w", err)
  }

  return db, nil
}

func dbMigration(db *gorm.DB) error {
  if err := db.AutoMigrate(&ItemMaster{}, &LatestItem{}); err != nil {
    return fmt.Errorf("db migration error: %w", err)
  }
  return nil
}
