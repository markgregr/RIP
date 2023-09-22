package api

type Baggage struct {
	BaggageID     int
	OwnerName     string
	Destination   string
	Weight        float64
	Status        string
	Src           string
	Href          string
	BaggageType   string
	Width, Height int
	Length        int
}