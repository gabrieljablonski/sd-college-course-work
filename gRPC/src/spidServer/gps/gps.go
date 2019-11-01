package gps

import (
	"fmt"
	"spidServer/requestHandling/protoBuffers"
)

type GlobalPosition struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

}

func NullPosition() GlobalPosition {
	return GlobalPosition{0, 0}
}

func FromProtoBufferEntity(position *protoBuffers.GlobalPosition) GlobalPosition {
	if position == nil {
		return NullPosition()
	}
	return GlobalPosition{
		Latitude:  position.Latitude,
		Longitude: position.Longitude,
	}
}

func (p GlobalPosition) String() string {
	return fmt.Sprintf("GlobalPosition{Latitude: %f, Longitude: %f}", p.Latitude, p.Longitude)
}

func (p GlobalPosition) ToProtoBufferEntity() *protoBuffers.GlobalPosition {
	return &protoBuffers.GlobalPosition{
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
	}
}