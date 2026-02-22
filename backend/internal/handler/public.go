package handler

import (
	"app/internal/models"
	"app/internal/service"

	"github.com/gofiber/fiber/v2"
)

type PublicHandler struct {
	msgSvc service.MessageService
}

func NewPublicHandler(msgSvc service.MessageService) *PublicHandler {
	return &PublicHandler{msgSvc: msgSvc}
}

func (h *PublicHandler) ListMessages(c *fiber.Ctx) error {
	msgs, err := h.msgSvc.ListPublicMessages(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if msgs == nil {
		msgs = []models.PublicMessage{}
	}
	return c.JSON(msgs)
}
