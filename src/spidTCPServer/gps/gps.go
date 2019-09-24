package gps

type GlobalPosition struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NullPosition() GlobalPosition {
	return GlobalPosition{0, 0}
}
