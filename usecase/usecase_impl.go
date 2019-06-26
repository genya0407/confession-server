package usecase

import (
	"github.com/genya0407/confession-server/entity"
)

// = func(AnonymousLoginInfoDTO, chatID, Socket)

type SocketImpl struct {
	c Chat
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

func GenerateJoinChatAnonymous(findChatByID entity.FindChatByID, registerAnonymousSocket entity.RegisterAnonymousSocket) JoinChatAnonymous {
	return func(anonLoginInfo AnonymousLoginInfoDTO, cID ChatID, s Socket) *JoinChatAnonymousError {
		chat, ok := findChatByID(cID)
		if !ok {
			s.Close()
			return &ChatNotFound
		}

		if chat.Anonymous.Token != anonLoginInfo.SessionToken {
			s.Close()
			return &InvalidToken
		}

		registerAnonymousSocket(cID, s)

		return nil
	}
}
