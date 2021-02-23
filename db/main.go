package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// Init inisializes the connection between the database and run the migration.
func Init() {
	var err error

	dsn := "catechdojo:#CATechDojo1017@tcp(127.0.0.1:3306)/catechdojo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("failed to establish the connection between the database", err)
	}

	err = db.AutoMigrate(&User{})

	if err != nil {
		log.Fatalln("failed to migrate the database", err)
	}
}
