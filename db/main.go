package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// Init inisializes the connection between the database and run the migration.
func Init() {
	var err error

	// TODO: 環境変数から読み取る
	dsn := os.Getenv("DB_ADDR")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("failed to establish the connection between the database", err)
	}

	log.Println("Connected to database")

	err = db.AutoMigrate(&User{}, &UserCharacter{}, &Character{})

	if err != nil {
		log.Fatalln("failed to migrate the database", err)
	}

	log.Println("Migration done")
}
