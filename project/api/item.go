package api

import (
	"time"
)

type Item struct {
	Entry       string `json:"entry" validate:"required,max=256"`
	OrderUid    string `json:"order_uid" validate:"required,ascii,max=256"`
	TrackNumber string `json:"track_number" validate:"required,ascii,max=14"`
	Delivery    struct {
		Name    string `json:"name" validate:"required,ascii,max=256"`
		Phone   string `json:"phone" validate:"required,e164"`
		Zip     string `json:"zip" validate:"required,number"`
		City    string `json:"city" validate:"required,max=256"`
		Address string `json:"address" validate:"required,max=256"`
		Region  string `json:"region" validate:"required,max=256"`
		Email   string `json:"email" validate:"required,email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string `json:"transaction" validate:"required,ascii,max=256"`
		RequestId    string `json:"request_id"`
		Currency     string `json:"currency" validate:"required,ascii,max=3"`
		Provider     string `json:"provider" validate:"required,ascii,max=256"`
		Amount       int    `json:"amount" validate:"required,number"`
		PaymentDt    int    `json:"payment_dt" validate:"required,number"`
		Bank         string `json:"bank" validate:"required,ascii,max=256"`
		DeliveryCost int    `json:"delivery_cost" validate:"required,number"`
		GoodsTotal   int    `json:"goods_total" validate:"required,number"`
		CustomFee    int    `json:"custom_fee" validate:"required,number"`
	} `json:"payment"`
	Items []struct {
		ChrtId      int    `json:"chrt_id" validate:"required,number"`
		TrackNumber string `json:"track_number" validate:"required,ascii,max=256"`
		Price       int    `json:"price" validate:"required,number"`
		Rid         string `json:"rid" validate:"required,ascii,max=256"`
		Name        string `json:"name" validate:"required,ascii,max=256"`
		Sale        int    `json:"sale" validate:"required,number"`
		Size        string `json:"size" validate:"required,number"`
		TotalPrice  int    `json:"total_price" validate:"required,number"`
		NmId        int    `json:"nm_id" validate:"required,number"`
		Brand       string `json:"brand" validate:"required,ascii,max=256"`
		Status      int    `json:"status" validate:"required,number"`
	} `json:"items"`
	Locale            string    `json:"locale" validate:"required,ascii,max=2"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id" validate:"required,ascii,max=256"`
	DeliveryService   string    `json:"delivery_service" validate:"required,ascii,max=256"`
	Shardkey          string    `json:"shardkey" validate:"required,number"`
	SmId              int       `json:"sm_id" validate:"required,number"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard" validate:"required,number"`
}
