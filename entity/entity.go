package entity

import (
	"sync"
	"time"

	"github.com/genya0407/confession-server/utils"

	"github.com/google/uuid"
)

// data structures

type TwitterAccount struct {
	TwitterAccountID string
	AccountID        uuid.UUID
}

type IAccount interface {
	AccountID() uuid.UUID
	Name() string
	ScreenName() string
	ImageURL() string
	Token() string
}

type Account struct {
	accountID  uuid.UUID
	name       string
	screenName string
	imageUrl   string
	token      string
}

func (a Account) AccountID() uuid.UUID {
	return a.accountID
}

func (a Account) Name() string {
	return a.name
}

func (a Account) ScreenName() string {
	return a.screenName
}

func (a Account) ImageURL() string {
	return a.imageUrl
}

func (a Account) Token() string {
	return a.token
}

type IAnonymous interface {
	Token() string
}

type Anonymous struct {
	token string
}

func (a Anonymous) Token() string {
	return a.token
}

const tokenLength = 30

func NewAnonymous() Anonymous {
	return Anonymous{token: utils.GenerateToken68Token(tokenLength)}
}

type IMessage interface {
	MessageID() uuid.UUID
	Text() TextMessage
	ByAnonymous() bool
	SentAt() time.Time
}

type Message struct {
	messageID   uuid.UUID
	text        TextMessage
	byAnonymous bool
	sentAt      time.Time
}

func (m *Message) MessageID() uuid.UUID {
	return m.messageID
}

func (m *Message) Text() TextMessage {
	return m.text
}

func (m *Message) ByAnonymous() bool {
	return m.byAnonymous
}

func (m *Message) SentAt() time.Time {
	return m.sentAt
}

func NewAnonymousMessage(text TextMessage) IMessage {
	return &Message{
		messageID:   utils.MustNewUUID(),
		text:        text,
		byAnonymous: true,
		sentAt:      time.Now(),
	}
}

func NewAccountMessage(text TextMessage) IMessage {
	return &Message{
		messageID:   utils.MustNewUUID(),
		text:        text,
		byAnonymous: false,
		sentAt:      time.Now(),
	}
}

type TextMessage = string

type Socket interface {
	SendText(msg IMessage)
	Close()
}

type ChatID = uuid.UUID

type IChat interface {
	ChatID() ChatID
	Messages() []IMessage
	StartedAt() time.Time
	FinishedAt() utils.NullableTime
	AuthorizeAnonymous(Anonymous) bool
	AuthorizeAccount(Account) bool
	SendAccountMessageToAnonymous(string)
	SendAnonymousMessageToAccount(string)
	RegisterAccountSocket(Socket)
	RegisterAnonymousSocket(Socket)
	Close()
}

type Chat struct {
	chatID          ChatID
	Account         Account
	Anonymous       Anonymous
	messages        []IMessage
	startedAt       time.Time
	finishedAt      utils.NullableTime
	AccountSocket   Socket
	AnonymousSocket Socket
	StoreChat       StoreChat
	m               *sync.Mutex
}

func NewChat(acc Account, anon Anonymous, storeChat StoreChat, beginningMessateText TextMessage) IChat {
	msgs := []IMessage{
		NewAnonymousMessage(beginningMessateText),
	}
	return &Chat{
		chatID:    utils.MustNewUUID(),
		Account:   acc,
		Anonymous: anon,
		messages:  msgs,
		startedAt: time.Now(),
		StoreChat: storeChat,
		m:         &sync.Mutex{},
	}
}

func (c *Chat) ChatID() uuid.UUID {
	return c.chatID
}

func (c *Chat) Messages() []IMessage {
	return c.messages
}

func (c *Chat) StartedAt() time.Time {
	return c.startedAt
}

func (c *Chat) FinishedAt() utils.NullableTime {
	return c.finishedAt
}

func (c *Chat) AuthorizeAccount(acc Account) bool {
	return c.Account == acc
}

func (c *Chat) AuthorizeAnonymous(anon Anonymous) bool {
	return c.Anonymous == anon
}

func (c *Chat) Close() {
	c.m.Lock()
	defer c.m.Unlock()

	c.finishedAt = utils.NullableTime{Null: false, Value: time.Now()}
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

	msg := NewAnonymousMessage(text)
	c.messages = append(c.messages, msg)
	c.StoreChat(c)

	if c.AccountSocket != nil {
		c.AccountSocket.SendText(msg)
	}
}

func (c *Chat) SendAccountMessageToAnonymous(text TextMessage) {
	c.m.Lock()
	defer c.m.Unlock()

	msg := NewAccountMessage(text)
	c.messages = append(c.messages, msg)
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

type FindAccountByID = func(uuid.UUID) (Account, bool)
type FindAccountByToken = func(string) (Account, bool)
type FindAnonymousByToken = func(string) (Anonymous, bool)

type FindChatByID = func(uuid.UUID) (Chat, bool)
type StoreChat = func(*Chat) error
