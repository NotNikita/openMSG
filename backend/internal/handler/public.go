package handler

import (
	"strconv"
	"time"

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
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	var before *time.Time
	if raw := c.Query("before"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid before timestamp, use RFC3339"})
		}
		before = &t
	}

	msgs, err := h.msgSvc.ListPublicMessages(c.Context(), limit, before)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if msgs == nil {
		msgs = []models.PublicMessage{}
	}
	return c.JSON(msgs)
}
