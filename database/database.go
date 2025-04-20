package database

import (
	"fmt"
	"log"

	"github.com/huuloc2026/go-to-do.git/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbCfg := config.AppConfig.DB

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
		dbCfg.Charset,
		dbCfg.ParseTime,
		dbCfg.Loc,
	)
	fmt.Println("DSN:", dsn)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to DB:", err)
	}

	log.Println("✅ Connected to MySQL DB!")
}
