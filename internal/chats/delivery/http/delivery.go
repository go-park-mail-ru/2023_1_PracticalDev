package http

import (
	"encoding/json"
	"fmt"
	pkgChats "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// UserID: Conn
var connectionsByUserID = map[int][]*websocket.Conn{}

// ChatID: Conn
var connectionsByChatID = map[int][]*websocket.Conn{}

const (
	chatsUrl        = "/chats"
	chatUrl         = "/chats/:id"
	chatMessagesUrl = "/chats/:id/messages"
	messagesUrl     = "/messages"
	wsChatUrl       = "/chat"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, serv pkgChats.Service) {
	del := delivery{serv, logger}
	// chats
	mux.GET(chatsUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.listByUser)), logger), logger))
	mux.GET(chatUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.get)), logger), logger))

	// messages
	mux.GET(chatMessagesUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.messagesList)), logger), logger))
	mux.GET(messagesUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.getMessagesByReceiver)), logger), logger))
	mux.GET(wsChatUrl, mw.HandleLogger(mw.ErrorHandler(authorizer(del.chatHandler), logger), logger))
}

type delivery struct {
	serv pkgChats.Service
	log  log.Logger
}

type NewChatMsg struct {
	Type string      `json:"type"`
	Chat models.Chat `json:"chat"`
}

func (del *delivery) NewChatSendToConns(chat models.Chat) {
	del.log.Debug("New chat for sending:", chat)

	newChatMsg := NewChatMsg{
		Type: "new_chat",
		Chat: chat,
	}

	// Необходимо разослать всем соединениям с user1ID и user2ID
	// по веб-сокету сообщение о том, что был создан новый чат
	// также добавить их соединения в список соединений чата
	conns, ok := connectionsByUserID[chat.User1ID]
	if ok {
		for _, con := range conns {
			con.WriteJSON(newChatMsg)
			connectionsByChatID[chat.ID] = append(connectionsByChatID[chat.ID], con)
		}
	} else {
		del.log.Debug("There are no connections for User1ID=", chat.User1ID)
	}

	conns, ok = connectionsByUserID[chat.User2ID]
	if ok {
		for _, con := range conns {
			con.WriteJSON(newChatMsg)
			connectionsByChatID[chat.ID] = append(connectionsByChatID[chat.ID], con)
		}
	} else {
		del.log.Debug("There are no connections for User2ID=", chat.User2ID)
	}

	del.log.Debug("New chat was sent to chat participants")
}

func (del *delivery) listByUser(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	chats, err := del.serv.ListByUser(userID)
	if err != nil {
		return err
	}

	response := newListResponse(chats)
	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

func (del *delivery) messagesList(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
	strID := p.ByName("id")
	chatID, err := strconv.Atoi(strID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidChatIDParam, err.Error())
	}

	messages, err := del.serv.MessagesList(chatID)
	if err != nil {
		return err
	}

	response := newMessagesListResponse(messages)
	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

func (del *delivery) getMessagesByReceiver(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	strReceiverID := r.URL.Query().Get("receiver_id")
	receiverID, err := strconv.Atoi(strReceiverID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	chat, err := del.serv.GetByUsers(userID, receiverID)
	if err != nil {
		return err
	}

	messages, err := del.serv.MessagesList(chat.ID)
	if err != nil {
		return err
	}

	response := newMessagesListResponse(messages)
	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

func (del *delivery) get(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
	strID := p.ByName("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidChatIDParam, err.Error())
	}

	chat, err := del.serv.Get(id)
	if err != nil {
		return err
	}

	response := newGetResponse(&chat)
	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

type NewMessageMsg struct {
	Type    string         `json:"type"`
	Message models.Message `json:"message"`
}

func (del *delivery) SendMessageToChatMembers(message models.Message) {
	del.log.Debug("New message for sending:", message)

	cons, ok := connectionsByChatID[message.ChatID]
	if !ok {
		del.log.Warn("There are no connections for this chat")
		return
	}

	newMessage := NewMessageMsg{
		Type:    "message",
		Message: message,
	}
	for _, con := range cons {
		con.WriteJSON(newMessage)
	}
	del.log.Debug("Message was sent to chat participants")
}

type msgRequest struct {
	Text       string `json:"text"`
	ReceiverID int    `json:"receiver_id"`
}

func (del *delivery) chatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}
	del.log.Debug("Handle new websocket request, userID =", userID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrUpgradeToWebSocket, err.Error())
	}
	defer func() {
		del.log.Debug(fmt.Sprintf("Close websocket connection=%p, userID=%d", conn, userID))
		err = conn.Close()
		if err != nil {
			del.log.Error(err)
		}
	}()
	del.log.Debug(fmt.Sprintf("Websocket connected: connection=%p, userID=%d", conn, userID))

	// Добавляем соединение в список соединений пользователя
	connectionsByUserID[userID] = append(connectionsByUserID[userID], conn)

	// Добавляем соединение в чаты
	chats, err := del.serv.ListByUser(userID)
	if err != nil {
		return err
	}
	for _, chat := range chats {
		connectionsByChatID[chat.ID] = append(connectionsByChatID[chat.ID], conn)
	}
	del.log.Debug(fmt.Sprintf("Сonnection %p was subscribed for chats events", conn))

	// Считываем сообщения пользователя
	for {
		del.log.Debug(fmt.Sprintf("Reading new message from connection %p...", conn))
		_, message, err := conn.ReadMessage()
		if err != nil {
			del.log.Error("error reading message from websocket: %v", err)
			break
		}

		msgReq := msgRequest{}
		err = json.Unmarshal(message, &msgReq)
		if err != nil {
			del.log.Error(err)
		}

		del.log.Debug(fmt.Sprintf("Got message (conn=%p; userID=%d): receiverID=%d; text=%s",
			conn, userID, msgReq.ReceiverID, msgReq.Text))

		var needChat bool
		chat, err := del.serv.GetByUsers(userID, msgReq.ReceiverID)
		if err != nil {
			if errors.Is(err, pkgErrors.ErrChatNotFound) {
				del.log.Debug(fmt.Sprintf("Not found chat for (user1ID=%d; user2ID=%d)", userID, msgReq.ReceiverID))

				// чата нет, нужно создать чат
				params := pkgChats.CreateParams{User1ID: userID, User2ID: msgReq.ReceiverID}
				createdChat, err := del.serv.Create(&params)
				if err != nil {
					return err
				}

				needChat = true

				chat = createdChat
			} else {
				return err
			}
		}

		// 1. Сохранить сообщение в БД
		params := pkgChats.SendMessageParams{AuthorID: userID, ChatID: chat.ID, Text: msgReq.Text}
		createdMessage, err := del.serv.SendMessage(&params)
		if err != nil {
			return err
		}

		if needChat {
			// Разослать всем соединениям, пользователи которых состоят в
			// новом чате сообщение о том, что был создан чат
			del.NewChatSendToConns(chat)
		}

		// 2. Разослать всем соединениям, пользователи которых состоят в данном чате,
		// созданное сообщение
		go del.SendMessageToChatMembers(*createdMessage)
	}

	return nil
}
