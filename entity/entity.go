package entity

import (
	"sync"
	"time"

	"github.com/genya0407/confession-server/utils"

	"github.com/google/uuid"
)

// utils

func mustNewUUID() uuid.UUID {
	u, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	return u
}

// data structures

type TwitterAccount struct {
	TwitterAccountID string
	AccountID        uuid.UUID
}

type Account struct {
	AccountID  uuid.UUID
	Name       string
	ScreenName string
	ImageUrl   string
	Token      string
}

type Anonymous struct {
	Token string
}

type Message struct {
	MessageID   uuid.UUID
	Text        TextMessage
	ByAnonymous bool
	SentAt      time.Time
}

type TextMessage = string

type Socket interface {
	SendText(msg Message)
	Close()
}

type ChatID = uuid.UUID

type IChat interface {
	ChatID() ChatID
	SendAccountMessageToAnonymous(string)
	SendAnonymousMessageToAccount(string)
	RegisterAccountSocket(Socket)
	RegisterAnonymousSocket(Socket)
	Close()
}

type FinishedAt struct {
	Finished bool
	t        time.Time
}

type Chat struct {
	chatID          ChatID
	Account         Account
	Anonymous       Anonymous
	Messages        []Message
	StartedAt       time.Time
	FinishedAt      FinishedAt
	AccountSocket   Socket
	AnonymousSocket Socket
	StoreChat       StoreChat
	m               *sync.Mutex
}

func NewChat(acc Account, anon Anonymous, storeChat StoreChat) IChat {
	return &Chat{
		chatID:    utils.MustNewUUID(),
		Account:   acc,
		Anonymous: anon,
		StartedAt: time.Now(),
		StoreChat: storeChat,
		m:         &sync.Mutex{},
	}
}

func (c *Chat) ChatID() uuid.UUID {
	return c.chatID
}

func (c *Chat) Close() {
	c.m.Lock()
	defer c.m.Unlock()

	c.FinishedAt = FinishedAt{true, time.Now()}
	c.StoreChat(c)

	if c.AccountSocket != nil {
		c.AccountSocket.Close()
	}

	if c.AnonymousSocket != nil {
		c.AnonymousSocket.Close()
	}
}

func (c *Chat) SendAnonymousMessageToAccount(text TextMessage) {
	c.m.Lock()
	defer c.m.Unlock()

	msg := Message{
		MessageID:   mustNewUUID(),
		Text:        text,
		ByAnonymous: true,
		SentAt:      time.Now(),
	}
	c.Messages = append(c.Messages, msg)
	c.StoreChat(c)

	if c.AccountSocket != nil {
		c.AccountSocket.SendText(msg)
	}
}

func (c *Chat) SendAccountMessageToAnonymous(text TextMessage) {
	c.m.Lock()
	defer c.m.Unlock()

	msg := Message{
		MessageID:   mustNewUUID(),
		Text:        text,
		ByAnonymous: false,
		SentAt:      time.Now(),
	}
	c.Messages = append(c.Messages, msg)
	c.StoreChat(c)

	if c.AnonymousSocket != nil {
		c.AnonymousSocket.SendText(msg)
	}
}

func (c *Chat) RegisterAccountSocket(s Socket) {
	c.m.Lock()
	defer c.m.Unlock()

	c.AccountSocket = s
}

func (c *Chat) RegisterAnonymousSocket(s Socket) {
	c.m.Lock()
	defer c.m.Unlock()

	c.AnonymousSocket = s
}

// repositories

type FindAccountByID = func(uuid.UUID) Account
type FindAccountByToken = func(string) Account
type FindAnonymousByToken = func(string) Anonymous

type FindChatByID = func(uuid.UUID) (Chat, bool)
type StoreChat = func(*Chat) error
