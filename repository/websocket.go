package repository

import (
	"sync"

	"github.com/genya0407/confession-server/domain"
)

type socketPair struct {
	accountSocket   domain.ISocket
	anonymousSocket domain.ISocket
}

type OnMemoryWebsocketStore struct {
	m       *sync.Mutex
	storage map[domain.ChatID]socketPair
}

func NewOnMemoryWebsocketStore() *OnMemoryWebsocketStore {
	return &OnMemoryWebsocketStore{
		m:       &sync.Mutex{},
		storage: map[domain.ChatID]socketPair{},
	}
}

func (wsm *OnMemoryWebsocketStore) FindAnonymousSocket(cID domain.ChatID) domain.ISocket {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair, ok := wsm.storage[cID]
	if !ok {
		return nil
	}

	return pair.anonymousSocket
}

func (wsm *OnMemoryWebsocketStore) FindAccountSocket(cID domain.ChatID) domain.ISocket {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair, ok := wsm.storage[cID]
	if !ok {
		return nil
	}

	return pair.accountSocket
}

func (wsm *OnMemoryWebsocketStore) RegisterAnonymousSocket(cID domain.ChatID, s domain.ISocket) {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair := wsm.storage[cID]
	pair.anonymousSocket = s
	wsm.storage[cID] = pair
}

func (wsm *OnMemoryWebsocketStore) RegisterAccountSocket(cID domain.ChatID, s domain.ISocket) {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair := wsm.storage[cID]
	pair.accountSocket = s
	wsm.storage[cID] = pair
}
