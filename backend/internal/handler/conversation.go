package handler

import (
	"app/internal/models"
	"app/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ConversationHandler struct {
	svc service.ConversationService
}

func NewConversationHandler(svc service.ConversationService) *ConversationHandler {
	return &ConversationHandler{svc: svc}
}

func (h *ConversationHandler) GetOrCreate(c *fiber.Ctx) error {
	var body struct {
		UserAID string `json:"user_a_id"`
		UserBID string `json:"user_b_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	conv, err := h.svc.GetOrCreateConversation(c.Context(), body.UserAID, body.UserBID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(conv)
}

func (h *ConversationHandler) GetByID(c *fiber.Ctx) error {
	conv, err := h.svc.GetConversation(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if conv == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "conversation not found"})
	}
	return c.JSON(conv)
}

func (h *ConversationHandler) ListByUser(c *fiber.Ctx) error {
	convs, err := h.svc.ListConversations(c.Context(), c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if convs == nil {
		convs = []models.Conversation{}
	}
	return c.JSON(convs)
}
