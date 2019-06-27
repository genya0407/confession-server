package domain

type ICreateNewChatService = func(IAccount, IAnonymous, MessageText) IChat

func GenerateCreateNewChatService(storeChat IStoreChat) ICreateNewChatService {
	return func(acc IAccount, anon IAnonymous, text MessageText) IChat {
		chat := NewChat(acc, anon, text)
		storeChat(chat)
		return chat
	}
}
