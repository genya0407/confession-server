package jsonapi

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/genya0407/confession-server/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type AccountJSON struct {
	AccountID  uuid.UUID `json:"account_id"`
	Name       string    `json:"name"`
	ScreenName string    `json:"screen_name"`
	ImageURL   string    `json:"image_url"`
}

type ChatJSON struct {
	ChatID     uuid.UUID     `json:"chat_id"`
	Account    AccountJSON   `json:"account"`
	Messages   []MessageJSON `json:"messages"`
	StartedAt  time.Time     `json:"started_at"` // Go uses ISO8601 format by default
	FinishedAt *time.Time    `json:"finished_at,omitempty"`
}

type NewMessageJSON struct {
	Text string `json:"text"`
}

type MessageJSON struct {
	MessageID   uuid.UUID `json:"message_id"`
	Text        string    `json:"text"`
	ByAnonymous bool      `json:"by_anonymous"`
	SentAt      time.Time `json:"sent_at"`
}

type ChatAbstractJSON struct {
	ChatID               uuid.UUID  `json:"chat_id"`
	BeginningMessageText string     `json:"beginning_message_text"`
	StartedAt            time.Time  `json:"started_at"`
	FinishedAt           *time.Time `json:"finished_at,omitempty"`
}

type CreateChatRequestJSON struct {
	BeginningMessageText string `json:"beginning_message_text"`
}

type CreateChatResponseJSON struct {
	ChatID                uuid.UUID `json:"chat_id"`
	AnonymousSessionToken string    `json:"session_token"`
}

type EverybodyEndpoint = func(http.ResponseWriter, *http.Request, httprouter.Params)
type AccountAuthorizedEndpoint = func(http.ResponseWriter, *http.Request, httprouter.Params, usecase.AccountLoginInfoDTO)
type AnonymousAuthorizedEndpoint = func(http.ResponseWriter, *http.Request, httprouter.Params, usecase.AnonymousLoginInfoDTO)

// func GenerateGetAccountInfo(getAccountInfo usecase.GetAccountInfo) EverybodyEndpoint {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		accountIDString := ps.ByName("account_id")
// 		accountID, err := uuid.Parse(accountIDString)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			fmt.Fprint(w, "Invalid AccountID")
// 			return
// 		}

// 		account, ok := getAccountInfo(accountID)
// 		if !ok {
// 			w.WriteHeader(http.StatusNotFound)
// 			fmt.Fprintf(w, "Account not found: %s", accountID.String())
// 			return
// 		}

// 		accountJSON := AccountJSON{
// 			AccountID:  account.AccountID,
// 			Name:       account.Name,
// 			ScreenName: account.ScreeName,
// 			ImageURL:   account.ImageURL,
// 		}
// 		json.NewEncoder(w).Encode(accountJSON)
// 	}
// }

// func GenerateCreateChat(createChat usecase.CreateChat) EverybodyEndpoint {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		accountIDString := ps.ByName("account_id")
// 		accountID, err := uuid.Parse(accountIDString)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			fmt.Fprint(w, "Invalid AccountID")
// 			return
// 		}

// 		createChatJSON := &CreateChatRequestJSON{}
// 		err = json.NewDecoder(r.Body).Decode(createChatJSON)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			fmt.Fprint(w, "Invalid request")
// 			return
// 		}

// 		createChatResultDTO, createChatErr := createChat(accountID, createChatJSON.BeginningMessageText)
// 		if createChatErr != nil {
// 			switch *createChatErr {
// 			case usecase.AccountNotFound:
// 				w.WriteHeader(http.StatusNotFound)
// 				fmt.Fprint(w, "Account not found")
// 			default:
// 				w.WriteHeader(http.StatusInternalServerError)
// 				fmt.Fprint(w, "Internal server error")
// 			}
// 			return
// 		}

// 		response := &CreateChatResponseJSON{
// 			ChatID:                createChatResultDTO.Chat.ChatID,
// 			AnonymousSessionToken: createChatResultDTO.AnonymousLoginInfo.SessionToken,
// 		}
// 		json.NewEncoder(w).Encode(response)
// 	}
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(_ *http.Request) bool { return true }, // TODO: probably insecure?
}

type SocketImpl struct {
	conn *websocket.Conn
}

func (s *SocketImpl) SendTextJSON(msg MessageJSON) {
	log.Printf(`SendTextJSON: %v`, msg)
	err := s.conn.WriteJSON(msg)
	if err != nil {
		log.Println(err.Error())
	}
}

func (s *SocketImpl) SendText(msg usecase.MessageDTO) {
	log.Printf(`SendText: %v`, msg)
	s.SendTextJSON(MessageJSON{
		MessageID:   msg.MessageID,
		Text:        msg.Text,
		ByAnonymous: msg.ByAnonymous,
		SentAt:      msg.SentAt,
	})
}

func (s *SocketImpl) Close() {
	s.conn.Close()
}

func GenerateJoinChatAnonymous(joinChatAnonymous usecase.JoinChatAnonymous, sendMessageAnonymousToAccount usecase.SendMessageAnonymousToAccount) AnonymousAuthorizedEndpoint {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, anonymousLoginInfo usecase.AnonymousLoginInfoDTO) {
		chatIDString := ps.ByName("chat_id")
		chatID, err := uuid.Parse(chatIDString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid ChatID")
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}

		joinChatAnonymous(anonymousLoginInfo, chatID, &SocketImpl{conn: conn})

		for {
			newMessageJSON := &NewMessageJSON{}
			err := conn.ReadJSON(newMessageJSON)
			log.Printf(`Message Received: %s`, newMessageJSON.Text)
			if err != nil {
				log.Println(err)
				conn.Close()
				break
			}
			err = sendMessageAnonymousToAccount(anonymousLoginInfo, chatID, newMessageJSON.Text)
			if err != nil {
				log.Println(err)
				conn.Close()
				break
			}
		}
	}
}

func GenerateJoinChatAccount(joinChatAccount usecase.JoinChatAccount, sendMessageAccountToAnonymous usecase.SendMessageAccountToAnonymous) AccountAuthorizedEndpoint {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, accountLoginInfo usecase.AccountLoginInfoDTO) {
		chatIDString := ps.ByName("chat_id")
		chatID, err := uuid.Parse(chatIDString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid ChatID")
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = joinChatAccount(accountLoginInfo, chatID, &SocketImpl{conn: conn})
		if err != nil {
			log.Println(err.Error())
			conn.Close()
			return
		}

		for {
			newMessageJSON := &NewMessageJSON{}
			err := conn.ReadJSON(newMessageJSON)
			log.Printf(`Message Received: %s`, newMessageJSON.Text)
			if err != nil {
				log.Println(err)
				conn.Close()
				break
			}
			err = sendMessageAccountToAnonymous(accountLoginInfo, chatID, newMessageJSON.Text)
			if err != nil {
				log.Println(err)
				conn.Close()
				break
			}
		}
	}
}
