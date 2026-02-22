package handler

import (
	"app/internal/models"
	"app/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	svc service.MessageService
}

func NewMessageHandler(svc service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) Send(c *fiber.Ctx) error {
	var body struct {
		ConversationID string `json:"conversation_id"`
		SenderID       string `json:"sender_id"`
		Ciphertext     string `json:"ciphertext"`
		Nonce          string `json:"nonce"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	msg, err := h.svc.SendMessage(c.Context(), body.ConversationID, body.SenderID, body.Ciphertext, body.Nonce)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(msg)
}

func (h *MessageHandler) ListByConversation(c *fiber.Ctx) error {
	msgs, err := h.svc.ListMessages(c.Context(), c.Params("conversationId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if msgs == nil {
		msgs = []models.Message{}
	}
	return c.JSON(msgs)
}
