package usecase

import (
	"errors"
	"github.com/genya0407/confession-server/entity"
)

// = func(AnonymousLoginInfoDTO, chatID, Socket)

type SocketImpl struct {
	s Socket
}

func (si *SocketImpl) SendText(msg entity.Message) {
	si.s.SendText(MessageDTO{
		MessageID:   msg.MessageID,
		Text:        msg.Text,
		SentAt:      msg.SentAt,
		ByAnonymous: msg.ByAnonymous,
	})
}

func (si *SocketImpl) Close() {
	si.s.Close()
}

func GenerateJoinChatAnonymous(findChatByID entity.FindChatAnonymous, registerAnonymousSocket entity.RegisterAnonymousSocket) JoinChatAnonymous {
	return func(anonLoginInfo AnonymousLoginInfoDTO, cID ChatID, s Socket) error {
		_, ok := findChatByID(cID, entity.Anonymous{Token: anonLoginInfo.SessionToken})
		if !ok {
			s.Close()
			return errors.New("Chat not found")
		}

		registerAnonymousSocket(cID, &SocketImpl{s: s})

		return nil
	}
}

func GenerateSendMessageAnonymousToAccount(findChatByID entity.FindChatAnonymous) SendMessageAnonymousToAccount {
	return func(anonLoginInfo AnonymousLoginInfoDTO, chatID ChatID, msgText MessageText) error {
		chat, ok := findChatByID(chatID, entity.Anonymous{Token: anonLoginInfo.SessionToken})
		if !ok {
			return errors.New("Chat not found")
		}

		chat.SendAnonymousMessageToAccount(msgText)

		return nil
	}
}
