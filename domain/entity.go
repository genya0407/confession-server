package domain

import (
	"sync"
	"time"

	"github.com/genya0407/confession-server/utils"

	"github.com/google/uuid"
)

// data structures

// type TwitterAccount struct {
// 	TwitterAccountID string
// 	AccountID        uuid.UUID
// }

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

func NewAccount(name string, screenName string, imageUrl string) IAccount {
	return Account{
		accountID:  utils.MustNewUUID(),
		name:       name,
		screenName: screenName,
		imageUrl:   imageUrl,
		token:      utils.GenerateToken68Token(tokenLength),
	}
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
	Text() MessageText
	ByAnonymous() bool
	SentAt() time.Time
}

type Message struct {
	messageID   uuid.UUID
	text        MessageText
	byAnonymous bool
	sentAt      time.Time
}

func (m Message) MessageID() uuid.UUID {
	return m.messageID
}

func (m Message) Text() MessageText {
	return m.text
}

func (m Message) ByAnonymous() bool {
	return m.byAnonymous
}

func (m Message) SentAt() time.Time {
	return m.sentAt
}

func NewAnonymousMessage(text MessageText) IMessage {
	return &Message{
		messageID:   utils.MustNewUUID(),
		text:        text,
		byAnonymous: true,
		sentAt:      time.Now(),
	}
}

func NewAccountMessage(text MessageText) IMessage {
	return &Message{
		messageID:   utils.MustNewUUID(),
		text:        text,
		byAnonymous: false,
		sentAt:      time.Now(),
	}
}

type MessageText = string

type ISocket interface {
	SendText(msg IMessage)
	Close()
}

type ChatID = uuid.UUID

type IChat interface {
	ChatID() ChatID
	Messages() []IMessage
	StartedAt() time.Time
	FinishedAt() utils.NullableTime
	Anonymous() IAnonymous
	Account() IAccount
	SendAccountMessageToAnonymous(string)
	SendAnonymousMessageToAccount(string)
	RegisterAccountSocket(ISocket)
	RegisterAnonymousSocket(ISocket)
	Close()
}

type Chat struct {
	chatID          ChatID
	account         IAccount
	anonymous       IAnonymous
	messages        []IMessage
	startedAt       time.Time
	finishedAt      utils.NullableTime
	accountSocket   ISocket
	anonymousSocket ISocket
	m               *sync.Mutex
}

func NewChat(acc IAccount, beginningMessateText MessageText) IChat {
	msgs := []IMessage{
		NewAnonymousMessage(beginningMessateText),
	}
	return &Chat{
		chatID:    utils.MustNewUUID(),
		account:   acc,
		anonymous: NewAnonymous(),
		messages:  msgs,
		startedAt: time.Now(),
		m:         &sync.Mutex{},
	}
}

func (c *Chat) ChatID() uuid.UUID {
	return c.chatID
}

func (c *Chat) Messages() []IMessage {
	return c.messages
}

func (c *Chat) Account() IAccount {
	return c.account
}

func (c *Chat) Anonymous() IAnonymous {
	return c.anonymous
}

func (c *Chat) StartedAt() time.Time {
	return c.startedAt
}

func (c *Chat) FinishedAt() utils.NullableTime {
	return c.finishedAt
}

func (c *Chat) Close() {
	c.m.Lock()
	defer c.m.Unlock()

	c.finishedAt = utils.NullableTime{Null: false, Value: time.Now()}

	if c.accountSocket != nil {
		c.accountSocket.Close()
	}

	if c.anonymousSocket != nil {
		c.anonymousSocket.Close()
	}
}

func (c *Chat) SendAnonymousMessageToAccount(text MessageText) {
	c.m.Lock()
	defer c.m.Unlock()

	msg := NewAnonymousMessage(text)
	c.messages = append(c.messages, msg)

	if c.accountSocket != nil {
		c.accountSocket.SendText(msg)
	}
}

func (c *Chat) SendAccountMessageToAnonymous(text MessageText) {
	c.m.Lock()
	defer c.m.Unlock()

	msg := NewAccountMessage(text)
	c.messages = append(c.messages, msg)

	if c.anonymousSocket != nil {
		c.anonymousSocket.SendText(msg)
	}
}

func (c *Chat) RegisterAccountSocket(s ISocket) {
	c.m.Lock()
	defer c.m.Unlock()

	c.accountSocket = s
}

func (c *Chat) RegisterAnonymousSocket(s ISocket) {
	c.m.Lock()
	defer c.m.Unlock()

	c.anonymousSocket = s
}
