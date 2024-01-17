package model

import "time"

type Delivery struct {
	DeliveryID     uint      `gorm:"type:serial;primarykey" json:"delivery_id"`
	FlightNumber   string    `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  *time.Time `json:"formation_date"`
	CompletionDate *time.Time `json:"completion_date"`
	DeliveryStatus string    `json:"delivery_status"`
	UserID         uint      `json:"user_id"`
	ModeratorID    *uint      `gorm:"foreignkey:delivery_id" json:"moderator_id"`
}

type DeliveryRequest struct {
	DeliveryID     uint      `json:"delivery_id"`
	FlightNumber   string    `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  *time.Time `json:"formation_date"`
	CompletionDate *time.Time `json:"completion_date"`
	DeliveryStatus string    `json:"delivery_status"`
	FullName 	   string 	 `json:"full_name"`
	ModeratorName string 	 `json:"moderator_name"`
}

type DeliveryGetResponse struct {
	DeliveryID     uint       `gorm:"foreignkey:delivery_id" json:"delivery_id"`
	FlightNumber   string     `json:"flight_number"`
	CreationDate   time.Time  `json:"creation_date"`
	FormationDate  *time.Time  `json:"formation_date"`
	CompletionDate *time.Time  `json:"completion_date"`
	DeliveryStatus string     `json:"delivery_status"`
	FullName       string     `json:"full_name"`
	Baggages       []Baggage  `gorm:"many2many:delivery_baggages" json:"baggages"`
}

type DeliveryUpdateFlightNumberRequest struct {
	FlightNumber   string    `json:"flight_number"`
}

type DeliveryUpdateStatusRequest struct {
	DeliveryStatus string    `json:"delivery_status"`
}