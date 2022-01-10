package db

import (
	"fmt"

	"gitlab.com/l3montree/cryptogotchi/clodhopper/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	User     string
	Password string
	Port     string
	DBName   string
	Host     string
}

func NewMySQL(config MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Port, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// automatic migrate all models
	// this will create all tables
	// and update all fields
	db.AutoMigrate(&models.Cryptogotchi{})
	db.AutoMigrate(&models.Record{})

	return db, nil
}
