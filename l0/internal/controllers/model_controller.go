package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/davoodmossgreen/wb/L0/internal/models"
	"github.com/jinzhu/gorm"

	"github.com/davoodmossgreen/wb/L0/internal/config"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

var db *gorm.DB

var (
	newModel models.Model
	newDelivery models.Delivery
	newItems models.Items
	newPayment models.Payment
	newOrder models.Order
)

var cache map[string]interface{} // кэш-мэп

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&models.Delivery{}, &models.Items{}, &models.Payment{}, &models.Order{})
	// Мигрируем столбцы структур в базу данных
}

func GetModel() models.Order {
	// Подключаемся к серверу nats
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Подключаемся к nats-streaming
	sc, err := stan.Connect("test-cluster", "client", stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, nats.DefaultURL)
	}

	// Получаем данные по каналу nats-streaming
	sub, err := sc.Subscribe("model",
  	func(m *stan.Msg) {
		json.Unmarshal(m.Data, &newModel)

		newDelivery = newModel.Delivery
		newItems = newModel.Items
		newPayment = newModel.Payment

		newOrder.Order_uid = newModel.Order_uid
		newOrder.Track_number = newModel.Track_number
		newOrder.Entry = newModel.Entry
		newOrder.Locale = newModel.Locale
		newOrder.Internal_signature = newModel.Internal_signature
		newOrder.Customer_id = newModel.Customer_id
		newOrder.Delivery_service = newModel.Delivery_service
		newOrder.Shardkey = newModel.Shardkey
		newOrder.Sm_id = newModel.Sm_id
		newOrder.Date_created = newModel.Date_created
		newOrder.Oof_shard = newModel.Oof_shard

		fmt.Println(newDelivery, newItems, newPayment, newOrder)
		cache[newOrder.Order_uid] = newOrder // Добавляем основные данные о заказе в кэш-мэп
	},
  	stan.StartWithLastReceived())

	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	defer sub.Close()
	return newOrder // Возвращаем основные данные о заказе
}

func Insert() {
		// Загружаем все данные о заказе из четырех структур в соответствующие четыре таблицы в БД
		db.NewRecord(newDelivery)
		db.Create(&newDelivery)

		db.NewRecord(newItems)
		db.Create(&newItems)

		db.NewRecord(newPayment)
		db.Create(&newPayment)

		db.NewRecord(newOrder)
		db.Create(&newOrder)
}


func GetOrderById(w http.ResponseWriter, r *http.Request) {
	var tpl *template.Template

	if r.Method == "POST" {
		uid := r.FormValue("order_uid") // Получаем от пользователя номер заказа

		// Проверяем соответствуют ли номеру заказа данные в кэш-мэп, если да - берем данные оттуда, если нет - проверяем наличие в БД
		if cache[uid] == nil {
			db.Raw("SELECT * FROM orders WHERE order_uid=?", uid).Scan(&newOrder)

			templ, err := tpl.ParseFiles("./templates/found.html")
			
			if err != nil {
				fmt.Println("error parsing template", err)
			}
					
			err = templ.Execute(w, newOrder)
			if err != nil {
				fmt.Println("error executing template(db)", err)
			}
		} else {
			templ, err := tpl.ParseFiles("./templates/found.html")
			
			if err != nil {
				fmt.Println("error parsing template", err)
			}
					
			err = templ.Execute(w, cache[uid])
			if err != nil {
				fmt.Println("error executing template(db)", err)
			}
		}
	} else { // При методе GET у пользователя высвечивается форма для заполнения номера заказа
		templ, err := tpl.ParseFiles("./templates/search.html")
	
		if err != nil {
				fmt.Println("error parsing template", err)
		}

		err = templ.Execute(w, nil)
		if err != nil {
			fmt.Println("error executing template search", err)
		}
	}
}