package gps

import (
	"fmt"
	"math"
	"spidServer/requestHandling/protoBuffers"
)

type GlobalPosition struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

}

func NullPosition() GlobalPosition {
	return GlobalPosition{math.NaN(), math.NaN()}
}

func (p GlobalPosition) IsValid() bool {
	return p == NullPosition() || (p.Latitude >= -90 && p.Latitude <= 90 && p.Longitude >= -180 && p.Longitude <= 180)
}

func FromProtoBufferEntity(position *protoBuffers.GlobalPosition) (GlobalPosition, error) {
	if position == nil {
		return NullPosition(), nil
	}
	p := GlobalPosition{
		Latitude:  position.Latitude,
		Longitude: position.Longitude,
	}
	if !p.IsValid() {
		return NullPosition(), fmt.Errorf("invalid latitude or longitude: %s", p)
	}
	return p, nil
}

func (p GlobalPosition) String() string {
	return fmt.Sprintf("GlobalPosition{Latitude: %f, Longitude: %f}", p.Latitude, p.Longitude)
}

func (p GlobalPosition) ToProtoBufferEntity() (*protoBuffers.GlobalPosition, error) {
	if !p.IsValid() {
		pbP, _ := NullPosition().ToProtoBufferEntity()
		return pbP, fmt.Errorf("invalid latitude or longitude: %s", p)
	}
	return &protoBuffers.GlobalPosition{
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
	}, nil
}