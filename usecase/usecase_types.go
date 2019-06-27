package usecase

import (
	"time"

	"github.com/genya0407/confession-server/utils"
	"github.com/google/uuid"
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
	FinishedAt utils.NullableTime
}

type MessageText = string

type MessageDTO struct {
	MessageID   uuid.UUID
	Text        MessageText
	ByAnonymous bool
	SentAt      time.Time
}

type AnonymousLoginInfoDTO struct {
	SessionToken string
}

type Socket interface {
	SendText(msg MessageDTO)
	Close()
}

type AccountID = uuid.UUID
type ChatID = uuid.UUID

// UseCases

/// Everybody

type GetAccountInfo = func(AccountID) (AccountInfoDTO, bool)
type GetFinishedChatAbstractsByAccountID = func(AccountID) []ChatAbstractDTO
type GetFinishedChatByChatID = func(ChatID) ChatDTO

/// Account

type RegisterAccount = func(TwitterAccountInfoDTO) AccountInfoDTO
type GetLoginAccountInfo = func(TwitterAccountInfoDTO) AccountLoginInfoDTO
type GetMyChatAbstracts = func(AccountLoginInfoDTO) []ChatAbstractDTO
type GetMyChat = func(AccountLoginInfoDTO, ChatID) ChatDTO
type SendMessageAccount = func(AccountLoginInfoDTO, ChatID, MessageDTO)
type FinishMyChat = func(AccountLoginInfoDTO, ChatID)

/// Anonymous

// CreateChat : When creating chat, one should specify beginning text
type CreateChatAndAnonymous = func(AccountID, string) (ChatDTO, AnonymousLoginInfoDTO, error)
type JoinChatAnonymous = func(AnonymousLoginInfoDTO, ChatID, Socket) error
type JoinChatAccount = func(AccountLoginInfoDTO, ChatID, Socket) error
type SendMessageAnonymousToAccount = func(AnonymousLoginInfoDTO, ChatID, MessageText) error

// UseCase Implementations
// TODO
