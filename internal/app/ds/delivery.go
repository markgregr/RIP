package ds

import "time"

type Delivery struct {
	DeliveryID     uint      `gorm:"autoIncrement;primarykey" json:"delivery_id"`
	FlightNumber   string    `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	DeliveryStatus string    `json:"delivery_status"`
	UserID         uint      `json:"user_id"`
	ModeratorID    uint      `json:"moderator_id"`
}
