package usecase

import (
	"github.com/google/uuid"
	"time"
)

// DTOs

type TwitterAccountInfoDTO struct {
	TwitterID  string
	ScreenName string
	Name       string
	ImageURL   string
}

type AccountInfoDTO struct {
	AccountID uuid.UUID
	Name      string
	ScreeName string
	ImageURL  string
}

type AccountLoginInfoDTO struct {
	SessionToken string
}

type ChatAbstractDTO struct {
	ChatID               uuid.UUID
	BeginningMessageText string
	StartedAt            time.Time
	FinishedAt           *time.Time
}

type ChatDTO struct {
	ChatID     uuid.UUID
	Messages   []MessageDTO
	StartedAt  time.Time
	FinishedAt *time.Time
}

type MessageDTO struct {
	MessageID   uuid.UUID
	Text        string
	ByAnonymous bool
	SentAt      time.Time
}

type AnonymousLoginInfoDTO struct {
	SessionToken string
}

type CreateChatResultDTO struct {
	Chat               ChatDTO
	AnonymousLoginInfo AnonymousLoginInfoDTO
}

// UseCases

type accountID = uuid.UUID
type chatID = uuid.UUID

/// Everybody

type GetAccountInfoUsecase = func(accountID) AccountInfoDTO
type GetFinishedChatAbstractsByAccountID = func(accountID) []ChatAbstractDTO
type GetFinishedChatByChatID = func(chatID) ChatDTO

/// Account

type RegisterAccount = func(TwitterAccountInfoDTO) AccountInfoDTO
type GetLoginAccountInfo = func(TwitterAccountInfoDTO) AccountLoginInfoDTO
type GetMyChatAbstracts = func(AccountLoginInfoDTO) []ChatAbstractDTO
type GetMyChat = func(AccountLoginInfoDTO, chatID) ChatDTO
type SendMessageAccount = func(AccountLoginInfoDTO, chatID, MessageDTO)
type FinishMyChat = func(AccountLoginInfoDTO, chatID)

/// Anonymous

// CreateChat : When creating chat, one should specify beginning text
type CreateChat = func(accountID, string) CreateChatResultDTO
type SendMessageAnonymous = func(AnonymousLoginInfoDTO, chatID, MessageDTO)

// UseCase Implementations
// TODO
