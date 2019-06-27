package domain

import "github.com/google/uuid"

type IFindAccountByID = func(uuid.UUID) (IAccount, bool)
type IFindAccountByToken = func(string) (IAccount, bool)
type IFindAnonymousByToken = func(string) (IAnonymous, bool)

type IFindChatByID = func(uuid.UUID) (IChat, bool)
type IStoreChat = func(IChat) error
