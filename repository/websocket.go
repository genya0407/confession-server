package repository

import (
	"sync"

	"github.com/genya0407/confession-server/entity"
)

type socketPair struct {
	accountSocket   entity.Socket
	anonymousSocket entity.Socket
}

type OnMemoryWebsocketStore struct {
	m       *sync.Mutex
	storage map[entity.ChatID]socketPair
}

func NewOnMemoryWebsocketStore() *OnMemoryWebsocketStore {
	return &OnMemoryWebsocketStore{
		m:       &sync.Mutex{},
		storage: map[entity.ChatID]socketPair{},
	}
}

func (wsm *OnMemoryWebsocketStore) FindAnonymousSocket(cID entity.ChatID) entity.Socket {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair, ok := wsm.storage[cID]
	if !ok {
		return nil
	}

	return pair.anonymousSocket
}

func (wsm *OnMemoryWebsocketStore) FindAccountSocket(cID entity.ChatID) entity.Socket {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair, ok := wsm.storage[cID]
	if !ok {
		return nil
	}

	return pair.accountSocket
}

func (wsm *OnMemoryWebsocketStore) RegisterAnonymousSocket(cID entity.ChatID, s entity.Socket) {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair := wsm.storage[cID]
	pair.anonymousSocket = s
	wsm.storage[cID] = pair
}

func (wsm *OnMemoryWebsocketStore) RegisterAccountSocket(cID entity.ChatID, s entity.Socket) {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair := wsm.storage[cID]
	pair.accountSocket = s
	wsm.storage[cID] = pair
}
