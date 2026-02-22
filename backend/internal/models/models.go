package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	PublicKey string    `json:"public_key"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

type Conversation struct {
	ID        string    `json:"id"`
	UserAID   string    `json:"user_a_id"`
	UserBID   string    `json:"user_b_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	Ciphertext     string    `json:"ciphertext"`
	Nonce          string    `json:"nonce"`
	CreatedAt      time.Time `json:"created_at"`
}

type PublicMessage struct {
	SenderNickname    string    `json:"sender_nickname"`
	RecipientNickname string    `json:"recipient_nickname"`
	Ciphertext        string    `json:"ciphertext"`
	CreatedAt         time.Time `json:"created_at"`
}
