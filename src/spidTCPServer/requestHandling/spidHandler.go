package requestHandling

import (
	"fmt"
	"github.com/google/uuid"
	"main/entities"
)

func (h Handler) querySpid(spidID string) (entities.Spid, error) {
	id, err := uuid.Parse(spidID)
	if err != nil {
		return entities.Spid{}, fmt.Errorf("invalid spid id: %s", err)
	}
	spid, err := h.Manager.QuerySpid(id)
	if err != nil {
		return entities.Spid{}, err
	}
	return spid, nil
}

func (h Handler) getSpidInfo(request Request) (response Response, ok bool) {
	return Response{}, false
}

func (h Handler) registerSpid(request Request) (response Response, ok bool) {
	return Response{}, false
}

func (h Handler) updateSpidLocation(request Request) (response Response, ok bool) {
	return Response{}, false
}

func (h Handler) deleteSpid(request Request) (response Response, ok bool) {
	return Response{}, false
}
