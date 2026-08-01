package database

import (
	"database/sql"
	"fmt"

	"cy11dsphk/server/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg config.DatabaseConfig) error {
	sqlDB, err := sql.Open("mysql", cfg.DSN(""))
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	createDatabaseSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		cfg.Name,
	)
	if _, err := sqlDB.Exec(createDatabaseSQL); err != nil {
		return err
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN(cfg.Name)), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}
