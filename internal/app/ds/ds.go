package ds

import "time"
const (
	BAGGAGE_STATUS_ACTIVE  = "активен"
	BAGGAGE_STATUS_DELETED = "удален"
)

type Baggage struct {
	BaggageID      int     `gorm:"primarykey" json:"baggage_id"`
	BaggageCode    string  `json:"baggage_code"`
	Weight         float32 `json:"weight"`
	Size           string  `json:"size"`
	BaggageStatus  string  `json:"baggage_status"`
	BaggageType    string  `json:"baggage_type"`
	OwnerName      string  `json:"owner_name"`
	PasportDetails string  `json:"pasport_details"`
	Airline        string  `json:"airline"`
	PhotoURL       string  `json:"photo"`
	Href           string  `json:"href"`
}
type Delivery struct {
	DeliveryID     int    `gorm:"primarykey" json:"delivery_id"`
	FlightNumber   string `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	DeliveryStatus string `json:"delivery_status"`
	UserID         int    `json:"user_id"`
	ModeratorID    int    `json:"moderator_id"`
}

type DeliveryBagggage struct {
	DeliveryBagggageID int `gorm:"primarykey" json:"delivery_baggage_id"`
	DeliveryID         int `json:"delivery_id"`
	BaggageID          int `json:"baggage_id"`
}

type User struct {
	UserID   int    `gorm:"primarykey" json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}