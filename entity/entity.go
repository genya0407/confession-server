package entity

import (
	"github.com/google/uuid"
	"time"
)

// data structures

type TwitterAccount struct {
	TwitterAccountID string
	AccountID        uuid.UUID
}

type Account struct {
	AccountID         uuid.UUID
	AccountName       string
	AccountScreenName string
	AccountImageUrl   string
	AccountToken      string
}

type Chat struct {
	Account    Account
	Anonymous  Anonymous
	Messages   []Message
	StartedAt  time.Time
	FinishedAt *time.Time
}

type Anonymous struct {
	AnonymousID    uuid.UUID
	AnonymousToken string
}

type Message struct {
	MessageID uuid.UUID
	Text      string
	SentAt    time.Time
}

// repositories

type FindAccountByID = func(uuid.UUID) Account
type FindAccountByToken = func(string) Account
type FindAnonymousByToken = func(string) Anonymous
type FindChatByID = func(uuid.UUID) Chat
