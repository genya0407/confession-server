package usecase

import (
	"errors"

	"github.com/genya0407/confession-server/domain"

	"github.com/google/uuid"
)

// utils

func messageToDTO(msg domain.IMessage) MessageDTO {
	return MessageDTO{
		MessageID:   msg.MessageID(),
		Text:        msg.Text(),
		SentAt:      msg.SentAt(),
		ByAnonymous: msg.ByAnonymous(),
	}
}

func chatToDTO(chat domain.IChat) ChatDTO {
	msgDTOs := []MessageDTO{}
	for _, msg := range chat.Messages() {
		msgDTOs = append(msgDTOs, messageToDTO(msg))
	}

	return ChatDTO{
		ChatID:     chat.ChatID(),
		Messages:   msgDTOs,
		StartedAt:  chat.StartedAt(),
		FinishedAt: chat.FinishedAt(),
	}
}

func authorizeAccount(chat domain.IChat, acc domain.IAccount) bool {
	return chat.Account().Token() == acc.Token()
}

func authorizeAnonymous(chat domain.IChat, anon domain.IAnonymous) bool {
	return chat.Anonymous().Token() == anon.Token()
}

// impls

type SocketImpl struct {
	s Socket
}

func (si *SocketImpl) SendText(msg domain.IMessage) {
	si.s.SendText(messageToDTO(msg))
}

func (si *SocketImpl) Close() {
	si.s.Close()
}

func GenerateCreateChatAnonymous(createNewChatService domain.ICreateNewChatService) CreateChatAndAnonymous {
	return func(accountID uuid.UUID, beginningMessageText string) (ChatDTO, AnonymousLoginInfoDTO, error) {
		chat, err := createNewChatService(accountID, beginningMessageText)
		if err != nil {
			return ChatDTO{}, AnonymousLoginInfoDTO{}, err
		}

		return chatToDTO(chat), AnonymousLoginInfoDTO{SessionToken: chat.Anonymous().Token()}, nil
	}
}

func GenerateJoinChatAnonymous(joinChatAnonymousService domain.IJoinChatAnonymousService, findAnonymousByToken domain.IFindAnonymousByToken, findChatByID domain.IFindChatByID) JoinChatAnonymous {
	return func(anonLoginInfo AnonymousLoginInfoDTO, cID ChatID, s Socket) error {
		anon, ok := findAnonymousByToken(anonLoginInfo.SessionToken)
		if !ok {
			return errors.New("Invalid token")
		}
		chat, ok := findChatByID(cID)
		if !ok || !authorizeAnonymous(chat, anon) { // authorization
			return errors.New("Chat not found")
		}

		joinChatAnonymousService(chat, &SocketImpl{s: s})
		return nil
	}
}

func GenerateJoinChatAccount(joinChatAccountService domain.IJoinChatAccountService, findAccountByToken domain.IFindAccountByToken, findChatByID domain.IFindChatByID) JoinChatAccount {
	return func(accLoginInfo AccountLoginInfoDTO, cID ChatID, s Socket) error {
		acc, ok := findAccountByToken(accLoginInfo.SessionToken)
		if !ok {
			return errors.New("Invalid token")
		}
		chat, ok := findChatByID(cID)
		if !ok || !authorizeAccount(chat, acc) { // authorization
			return errors.New("Chat not found")
		}

		joinChatAccountService(chat, &SocketImpl{s: s})
		return nil
	}
}

func GenerateSendMessageAnonymousToAccount(
	sendAnonymousMessageToAccountService domain.ISendAnonymousMessageToAccountService,
	findAnonymousByToken domain.IFindAnonymousByToken,
	findChatByID domain.IFindChatByID,
) SendMessageAnonymousToAccount {
	return func(anonLoginInfo AnonymousLoginInfoDTO, chatID ChatID, text MessageText) error {
		anon, ok := findAnonymousByToken(anonLoginInfo.SessionToken)
		if !ok {
			return errors.New("Invalid token")
		}
		chat, ok := findChatByID(chatID)
		if !ok || !authorizeAnonymous(chat, anon) { // authorization
			return errors.New("Chat not found")
		}

		sendAnonymousMessageToAccountService(chat, text)
		return nil
	}
}

func GenerateSendMessageAccountToAnonymous(
	sendAccountMessageToAnonymousService domain.ISendAccountMessageToAnonymousService,
	findAccountByToken domain.IFindAccountByToken,
	findChatByID domain.IFindChatByID,
) SendMessageAccountToAnonymous {
	return func(accLoginInfo AccountLoginInfoDTO, chatID ChatID, text MessageText) error {
		acc, ok := findAccountByToken(accLoginInfo.SessionToken)
		if !ok {
			return errors.New("Invalid token")
		}
		chat, ok := findChatByID(chatID)
		if !ok || !authorizeAccount(chat, acc) { // authorization
			return errors.New("Chat not found")
		}

		sendAccountMessageToAnonymousService(chat, text)
		return nil
	}
}
