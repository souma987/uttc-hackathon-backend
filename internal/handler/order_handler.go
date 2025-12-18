package handler

import "uttc-hackathon-backend/internal/service"

type OrderHandler struct {
	svc     *service.OrderService
	userSvc *service.UserService
}

func NewOrderHandler(svc *service.OrderService, userSvc *service.UserService) *OrderHandler {
	return &OrderHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}
