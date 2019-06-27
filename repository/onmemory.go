package repository

import (
	"sync"

	"github.com/genya0407/confession-server/entity"
	"github.com/google/uuid"
)

type OnMemoryChatRepository struct {
	websocketStore *OnMemoryWebsocketStore
	chatStore      map[uuid.UUID]entity.IChat
	m              *sync.Mutex
}

func NewOnMemoryChatRepository() *OnMemoryChatRepository {
	return &OnMemoryChatRepository{
		websocketStore: NewOnMemoryWebsocketStore(),
		chatStore:      map[uuid.UUID]entity.IChat{},
		m:              &sync.Mutex{},
	}
}

func (repo *OnMemoryChatRepository) Store(chat entity.IChat) error {
	repo.chatStore[chat.ChatID()] = chat
	return nil
}

func (repo *OnMemoryChatRepository) FindByID(chatID uuid.UUID) (entity.IChat, bool) {
	chat, ok := repo.chatStore[chatID]
	return chat, ok
}
