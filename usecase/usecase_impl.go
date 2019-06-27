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
		_, ok := findAnonymousByToken(anonLoginInfo.SessionToken)
		if !ok {
			return errors.New("Invalid token")
		}
		chat, ok := findChatByID(cID)
		if !ok || chat.Anonymous().Token() != anonLoginInfo.SessionToken { // authorization
			return errors.New("Chat not found")
		}

		joinChatAnonymousService(chat, &SocketImpl{s: s})
		return nil
	}
}

func GenerateJoinChatAccount(joinChatAccountService domain.IJoinChatAccountService, findAccountByToken domain.IFindAccountByToken, findChatByID domain.IFindChatByID) JoinChatAccount {
	return func(accLoginInfo AccountLoginInfoDTO, cID ChatID, s Socket) error {
		_, ok := findAccountByToken(accLoginInfo.SessionToken)
		if !ok {
			return errors.New("Invalid token")
		}
		chat, ok := findChatByID(cID)
		if !ok || chat.Account().Token() != accLoginInfo.SessionToken { // authorization
			return errors.New("Chat not found")
		}

		joinChatAccountService(chat, &SocketImpl{s: s})
		return nil
	}
}

// func GenerateSendMessageAnonymousToAccount(findChatByID domain.FindChatAnonymous) SendMessageAnonymousToAccount {
// 	return func(anonLoginInfo AnonymousLoginInfoDTO, chatID ChatID, msgText MessageText) error {
// 		chat, ok := findChatByID(chatID, domain.Anonymous{Token: anonLoginInfo.SessionToken})
// 		if !ok {
// 			return errors.New("Chat not found")
// 		}

// 		chat.SendAnonymousMessageToAccount(msgText)

// 		return nil
// 	}
// }
