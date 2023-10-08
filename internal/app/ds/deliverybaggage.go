package ds

type DeliveryBaggage struct {
	DeliveryID uint `gorm:"autoIncrement;primaryKey;index" json:"delivery_id"`
	BaggageID  uint `gorm:"autoIncrement;primaryKey;index" json:"baggage_id"`
}