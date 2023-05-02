package http

import (
	"encoding/json"
	"fmt"
	pkgChats "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
	ws "github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// UserID: Conn
var connectionsByUserID = map[int][]*ws.Conn{}

const (
	chatsUrl = "/chats"
	chatUrl  = "/chats/:id"

	chatMessagesUrl = "/chats/:id/messages"
	messagesUrl     = "/messages"
	wsChatUrl       = "/chat"
)

type delivery struct {
	serv     pkgChats.Service
	log      log.Logger
	upgrader ws.Upgrader
}

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, csrf mw.CSRFMiddleware, serv pkgChats.Service) {
	del := delivery{
		serv,
		logger,
		ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}}
	// chats
	mux.GET(chatsUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.listByUser))), logger), logger))
	mux.GET(chatUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.get))), logger), logger))

	// messages
	mux.GET(chatMessagesUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.messagesList))), logger), logger))
	mux.GET(messagesUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.getMessagesByReceiver))), logger), logger))

	// connect to websocket
	mux.GET(wsChatUrl, mw.HandleLogger(mw.ErrorHandler(authorizer(del.chatHandler), logger), logger))
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

func (del *delivery) sendResponseToUser(response any, userID int) error {
	del.log.Debug(fmt.Sprintf("New response %v for sending to userID=%d:", response, userID))
	conns, ok := connectionsByUserID[userID]
	if ok {
		for _, con := range conns {
			err := con.WriteJSON(response)
			if err != nil {
				del.log.Error("sendResponseToUser: error: %v", err)
				return err
			}
		}
	} else {
		del.log.Debug("There are no connections for userID=", userID)
	}

	return nil
}

func (del *delivery) sendNewChatToChatMembers(chat models.Chat) error {
	del.log.Debug("New chat for sending:", chat)
	newChat := newChatResponse{
		Type: "new_chat",
		Chat: chat,
	}

	// Необходимо разослать всем соединениям с User1ID и User2ID
	// по веб-сокету сообщение о том, что был создан новый чат
	err := del.sendResponseToUser(newChat, chat.User1ID)
	if err != nil {
		return err
	}
	err = del.sendResponseToUser(newChat, chat.User2ID)
	if err != nil {
		return err
	}

	del.log.Debug("New chat was sent to chat participants")
	return nil
}

func (del *delivery) sendMessageToChatMembers(message models.Message, user1ID, user2ID int) error {
	del.log.Debug("New message for sending:", message)
	newMessage := newMessageResponse{
		Type:    "message",
		Message: message,
	}

	// Необходимо разослать всем соединениям с user1ID и user2ID
	// по веб-сокету сообщение о том, что было создано новое сообщение
	err := del.sendResponseToUser(newMessage, user1ID)
	if err != nil {
		return err
	}
	err = del.sendResponseToUser(newMessage, user2ID)
	if err != nil {
		return err
	}

	del.log.Debug("New message was sent to chat participants")
	return nil
}

func (del *delivery) chatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}
	del.log.Debug("Handle new websocket request, userID =", userID)

	conn, err := del.upgrader.Upgrade(w, r, nil)
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

	// Считываем сообщения пользователя
	for {
		del.log.Debug(fmt.Sprintf("Start reading new messages from connection %p...", conn))
		_, message, err := conn.ReadMessage()
		if err != nil {
			del.log.Error("error reading message from websocket: %v", err)
			return err
		}

		msgReq := msgRequest{}
		err = json.Unmarshal(message, &msgReq)
		if err != nil {
			del.log.Error(err)
			return err
		}
		del.log.Debug(fmt.Sprintf("Got message (conn=%p; userID=%d): receiverID=%d; text=%s",
			conn, userID, msgReq.ReceiverID, msgReq.Text))

		var sendNewChat bool
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
				sendNewChat = true
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

		// 2. Рассылаем новый чат только при успешной записи первого сообщения в БД
		if sendNewChat {
			// Разослать всем соединениям, пользователи которых состоят в
			// новом чате, сообщение о том, что был создан чат
			err = del.sendNewChatToChatMembers(chat)
			if err != nil {
				return err
			}
		}

		// 3. Разослать всем соединениям, пользователи которых состоят в данном чате,
		// созданное сообщение
		go func() {
			err := del.sendMessageToChatMembers(*createdMessage, chat.User1ID, chat.User2ID)
			if err != nil {
				del.log.Error(err)
			}
		}()
	}

	return nil
}
