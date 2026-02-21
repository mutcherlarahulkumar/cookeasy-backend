package handlers

import "cookeasy/domain"

type Handler struct {
	service domain.Service
}

func NewHandler(s domain.Service) *Handler {
	return &Handler{
		service: s,
	}
}
