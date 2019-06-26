package repositoryimpl

import (
	"github.com/genya0407/confession-server/entity"
	"sync"
)

type socketPair struct {
	accountSocket   entity.Socket
	anonymousSocket entity.Socket
}

type WebsocketStoreOnMemory struct {
	m       sync.Mutex
	storage map[entity.ChatID]socketPair
}

func (wsm *WebsocketStoreOnMemory) FindAnonymousSocket(cID entity.ChatID) entity.Socket {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair, ok := wsm.storage[cID]
	if !ok {
		return nil
	}

	return pair.anonymousSocket
}

func (wsm *WebsocketStoreOnMemory) FindAccountSocket(cID entity.ChatID) entity.Socket {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair, ok := wsm.storage[cID]
	if !ok {
		return nil
	}

	return pair.accountSocket
}

func (wsm *WebsocketStoreOnMemory) RegisterAnonymousSocket(cID entity.ChatID, s entity.Socket) {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair := wsm.storage[cID]
	pair.anonymousSocket = s
	wsm.storage[cID] = pair
}

func (wsm *WebsocketStoreOnMemory) RegisterAccountSocket(cID entity.ChatID, s entity.Socket) {
	wsm.m.Lock()
	defer wsm.m.Unlock()

	pair := wsm.storage[cID]
	pair.accountSocket = s
	wsm.storage[cID] = pair
}
