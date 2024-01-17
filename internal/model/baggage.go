// Package model ...
package model

// Baggage представляет информацию о багаже.
type Baggage struct {
	BaggageID      uint    `gorm:"type:serial;primarykey" json:"baggage_id"`
	BaggageCode    string  `json:"baggage_code" example:"ABC123"`
	Weight         float32 `json:"weight" example:"23.5"`
	Size           string  `json:"size" example:"large"`
	BaggageStatus  string  `json:"baggage_status" example:"checked"`
	BaggageType    string  `json:"baggage_type" example:"suitcase"`
	OwnerName      string  `json:"owner_name" example:"John Doe"`
	PasportDetails string  `json:"pasport_details" example:"123456789"`
	Airline        string  `json:"airline" example:"AirlineX"`
	PhotoURL       string  `json:"photo_url" example:"http://example.com/baggage.jpg"`
}

// BaggageRequest представляет запрос на создание багажа.
type BaggageRequest struct {
	BaggageCode    string  `json:"baggage_code" example:"ABC123"`
	Weight         float32 `json:"weight" example:"23.5"`
	Size           string  `json:"size" example:"large"`
	BaggageType    string  `json:"baggage_type" example:"suitcase"`
	OwnerName      string  `json:"owner_name" example:"John Doe"`
	PasportDetails string  `json:"pasport_details" example:"123456789"`
	Airline        string  `json:"airline" example:"AirlineX"`
}

// BaggagesGetResponse представляет ответ с информацией о багажах и идентификаторе доставки.
type BaggagesGetResponse struct {
	Baggages   []Baggage `json:"baggages"`
	DeliveryID uint      `json:"delivery_id" example:"1"`
}
