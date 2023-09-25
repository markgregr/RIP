package ds

const (
	BAGGAGE_STATUS_ACTIVE  = "активен"
	BAGGAGE_STATUS_DELETED = "удален"
)

type Baggage struct {
	BaggageID     int     `gorm:"primarykey" json:"baggage_id"`
	BaggageCode   string  `json:"baggage_code"`
	Weight        float32 `json:"weight"`
	Size          string  `json:"size"`
	BaggageStatus string  `json:"baggage_status"`
	BaggageType   string  `json:"baggage_type"`
	Destination   string  `json:"destination"`
	PhotoURL      string  `json:"photo"`
	Href          string  `json:"href"`
}
type Delivery struct {
	DeliveryID       int    `gorm:"primarykey" json:"delivery_id"`
	FlightNumber     string `json:"flight_number"`
	RegistrationDate string `json:"registration_date"`
	LoadingDate      string `json:"loading_date"`
	ClaimDate        string `json:"claim_date"`
	DeliveryStatus   string `json:"delivery_status"`
	PassengerID      int    `json:"passenger_id"`
	ModeratorID      int    `json:"moderator_id"`
}

type DeliveryBagggage struct {
	DeliveryBagggageID int `gorm:"primarykey" json:"delivery_baggage_id"`
	DeliveryID         int `json:"delivery_id"`
	BaggageID          int `json:"baggage_id"`
}

type Passenger struct {
	PassengerID     int    `gorm:"primarykey" json:"passenger_id"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	PassportDetails string `json:"passport_details"`
}