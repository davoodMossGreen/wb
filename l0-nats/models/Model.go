package models

import (
	"time"

)

type Model struct {
	Order_uid string `json:"order_uid"`
	Track_number string `json:"track_number"`
	Entry string `json:"entry"`
	Delivery Delivery `json:"delivery" gorm:"embedded;embeddedPrefix:delivery_"`
	Payment Payment `json:"payment" gorm:"embedded;embeddedPrefix:payment_"`
	Items Items `json:"items" gorm:"embedded;embeddedPrefix:items_"`
	Locale string `json:"locale"`
	Internal_signature string `json:"internal_signature"`
	Customer_id string `json:"customer_id"`
	Delivery_service string `json:"delivery_service"`
	Shardkey string `json:"shardkey"`
	Sm_id int `json:"sm_id"`
	Date_created time.Time `json:"date_created"`
	Oof_shard string `json:"oof_shard"`
}