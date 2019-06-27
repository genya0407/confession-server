package domain

import "github.com/google/uuid"

type IFindAccountByID = func(uuid.UUID) (Account, bool)
type IFindAccountByToken = func(string) (Account, bool)
type IFindAnonymousByToken = func(string) (Anonymous, bool)

type IFindChatByID = func(uuid.UUID) (Chat, bool)
type IStoreChat = func(IChat) error
