package ds

type Baggage struct {
	BaggageID      uint    `gorm:"type:serial;primarykey" json:"baggage_id"`
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
