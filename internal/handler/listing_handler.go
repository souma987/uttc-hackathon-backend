package handler

import (
	"uttc-hackathon-backend/internal/service"
)

type ListingHandler struct {
	svc     *service.ListingService
	userSvc *service.UserService
}

func NewListingHandler(svc *service.ListingService, userSvc *service.UserService) *ListingHandler {
	return &ListingHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}
