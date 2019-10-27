package gps

import "spidServer/requestHandling/grpcWrapper/spidProtoBuffers"

type GlobalPosition struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

}

func NullPosition() GlobalPosition {
	return GlobalPosition{0, 0}
}

func FromProtoBufferEntity(position *spidProtoBuffers.GlobalPosition) GlobalPosition {
	return GlobalPosition{
		Latitude:  position.Latitude,
		Longitude: position.Longitude,
	}
}

func (p GlobalPosition) ToProtoBufferEntity() *spidProtoBuffers.GlobalPosition {
	return &spidProtoBuffers.GlobalPosition{
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
	}
}