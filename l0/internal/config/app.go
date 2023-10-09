package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

	

func Connect() {
	// Получаем данные из .env
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbport := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// Подключаемся к базе данных
	dbURI:= fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", host, user, name, password, dbport)
	
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Подключение к базе данных прошло успешно")
	}
}

func GetDB() *gorm.DB {
	return db
}
