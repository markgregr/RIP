package ds

type DeliveryBaggage struct {
	DeliveryID uint `gorm:"type:serial;primaryKey;index" json:"delivery_id"`
	BaggageID  uint `gorm:"type:serial;primaryKey;index" json:"baggage_id"`
}