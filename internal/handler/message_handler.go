package handler

import (
	"uttc-hackathon-backend/internal/service"
)

type MessageHandler struct {
	svc     *service.MessageService
	userSvc *service.UserService
}

func NewMessageHandler(svc *service.MessageService, userSvc *service.UserService) *MessageHandler {
	return &MessageHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}
