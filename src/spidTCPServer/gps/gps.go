package gps

type GlobalPosition struct {
	Latitude float64
	Longitude float64
}

func NullPosition() GlobalPosition {
	return GlobalPosition{0, 0}
}
