package repository

import (
	"github.com/k0kubun/pp"
	"log"
	"sync"

	"github.com/genya0407/confession-server/domain"
	"github.com/google/uuid"
)

type OnMemoryRepository struct {
	ChatStorage      map[uuid.UUID]domain.IChat
	AnonymousStorage map[string]domain.IAnonymous
	AccountStorage   map[string]domain.IAccount
	m                *sync.Mutex
}

func NewOnMemoryRepository() *OnMemoryRepository {
	return &OnMemoryRepository{
		ChatStorage:      map[uuid.UUID]domain.IChat{},
		AnonymousStorage: map[string]domain.IAnonymous{},
		AccountStorage:   map[string]domain.IAccount{},
		m:                &sync.Mutex{},
	}
}

func (repo *OnMemoryRepository) StoreChat(chat domain.IChat) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	repo.ChatStorage[chat.ChatID()] = chat
	repo.AnonymousStorage[chat.Anonymous().Token()] = chat.Anonymous()
	return nil
}

func (repo *OnMemoryRepository) FindChatByID(chatID uuid.UUID) (domain.IChat, bool) {
	repo.m.Lock()
	defer repo.m.Unlock()

	chat, ok := repo.ChatStorage[chatID]
	return chat, ok
}

func (repo *OnMemoryRepository) FindAccountByToken(token string) (domain.IAccount, bool) {
	repo.m.Lock()
	defer repo.m.Unlock()

	acc, ok := repo.AccountStorage[token]
	log.Println(`FindAccountByToken tried.`)
	if !ok {
		log.Println(`But failed (not found)`)
	}
	pp.Print(acc)
	return acc, ok
}

func (repo *OnMemoryRepository) FindAnonymousByToken(token string) (domain.IAnonymous, bool) {
	repo.m.Lock()
	defer repo.m.Unlock()

	acc, ok := repo.AnonymousStorage[token]
	return acc, ok
}
