package main

import (

	"github.com/davoodmossgreen/wb/L0/internal/controllers"

	
	"log"
	"net/http"
	"github.com/gorilla/mux"

)



func main() {
	r := mux.NewRouter() // Иницируем роутер
	r.HandleFunc("/", controllers.GetOrderById).Methods("GET", "POST") // Привязываем роутер к функции GetOrderById
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9010", r)) // Запускаем локальный сервер
}


