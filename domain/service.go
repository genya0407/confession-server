package domain

import (
	"errors"

	"github.com/google/uuid"
)

type ICreateNewChatService = func(uuid.UUID, MessageText) (IChat, error)

func GenerateCreateNewChatService(storeChat IStoreChat, findAccountByID IFindAccountByID) ICreateNewChatService {
	return func(accID uuid.UUID, text MessageText) (IChat, error) {
		acc, ok := findAccountByID(accID)
		if !ok {
			return &Chat{}, errors.New("Account not found")
		}

		chat := NewChat(acc, text)
		storeChat(chat)

		return chat, nil
	}
}

type IJoinChatAnonymousService = func(IChat, ISocket)

func GenerateJoinChatAnonymousService(findChatByID IFindChatByID, storeChat IStoreChat) IJoinChatAnonymousService {
	return func(chat IChat, s ISocket) {
		chat.RegisterAnonymousSocket(s)
		storeChat(chat)
	}
}

type IJoinChatAccountService = func(IChat, ISocket)

func GenerateJoinChatAccountService(findChatByID IFindChatByID, storeChat IStoreChat) IJoinChatAnonymousService {
	return func(chat IChat, s ISocket) {
		chat.RegisterAccountSocket(s)
		storeChat(chat)
	}
}

type ISendAccountMessageToAnonymousService = func(IChat, MessageText)

func GenerateSendAccountMessageToAnonymousService(storeChat IStoreChat) ISendAccountMessageToAnonymousService {
	return func(chat IChat, text MessageText) {
		chat.SendAccountMessageToAnonymous(text)
		storeChat(chat)
	}
}

type ISendAnonymousMessageToAccountService = func(IChat, MessageText)

func GenerateSendAnonymousMessageToAnonymousService(storeChat IStoreChat) ISendAnonymousMessageToAccountService {
	return func(chat IChat, text MessageText) {
		chat.SendAnonymousMessageToAccount(text)
		storeChat(chat)
	}
}

type ICloseChat = func(IChat)

func GenerateCloseChatService(storeChat IStoreChat) ICloseChat {
	return func(chat IChat) {
		chat.Close()
		storeChat(chat)
	}
}
