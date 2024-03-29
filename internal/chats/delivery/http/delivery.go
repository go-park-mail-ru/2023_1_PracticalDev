package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/connectionservice"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/xss"
	"net/http"
	"strconv"

	ws "github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	pkgChats "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

const (
	chatsUrl = "/chats"
	chatUrl  = "/chats/:id"

	chatMessagesUrl = "/chats/:id/messages"
	messagesUrl     = "/messages"
	wsChatUrl       = "/chat"
)

type delivery struct {
	serv        pkgChats.Service
	connService *connectionservice.Service
	log         *zap.Logger
	upgrader    ws.Upgrader
}

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger, authorizer mw.Authorizer, csrf mw.CSRFMiddleware, serv pkgChats.Service) {
	del := delivery{
		serv:        serv,
		connService: connectionservice.NewService(logger),
		log:         logger,
		upgrader: ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}}

	// chats
	mux.GET(chatsUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.ListByUser))), logger), logger))
	mux.GET(chatUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.Get))), logger), logger))

	// messages
	mux.GET(chatMessagesUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.MessagesList))), logger), logger))
	mux.GET(messagesUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(csrf(del.GetMessagesByReceiver))), logger), logger))

	// connect to websocket
	mux.GET(wsChatUrl, mw.HandleLogger(mw.ErrorHandler(authorizer(del.chatHandler), logger), logger))
}

func (del *delivery) ListByUser(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
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
	data, err := response.MarshalJSON()
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

func (del *delivery) MessagesList(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
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
	data, err := response.MarshalJSON()
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

func (del *delivery) GetMessagesByReceiver(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
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
	data, err := response.MarshalJSON()
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

func (del *delivery) Get(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
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
	data, err := response.MarshalJSON()
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

func (del *delivery) sendNewChatToChatMembers(chat models.Chat) error {
	del.log.Debug("New chat for sending", zap.Any("chat", chat))
	newChat := newChatResponse{
		Type: "new_chat",
		Chat: chat,
	}

	err := del.connService.Broadcast(newChat, chat.User1ID)
	if err != nil {
		return err
	}
	err = del.connService.Broadcast(newChat, chat.User2ID)
	if err != nil {
		return err
	}

	del.log.Debug("New chat was sent to chat participants")
	return nil
}

func (del *delivery) sendMessageToChatMembers(message models.Message, user1ID, user2ID int) error {
	del.log.Debug("New message for sending", zap.Any("message", message))

	message.Text = xss.Sanitize(message.Text)
	newMessage := newMessageResponse{
		Type:    "message",
		Message: message,
	}

	err := del.connService.Broadcast(newMessage, user1ID)
	if err != nil {
		return err
	}
	err = del.connService.Broadcast(newMessage, user2ID)
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
	del.log.Debug("Handle new websocket request", zap.Int("user_id", userID))

	conn, err := del.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrUpgradeToWebSocket, err.Error())
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			del.log.Error(constants.FailedCloseConnection, zap.Error(err))
		}
		del.log.Debug("Websocket closed", zap.Int("user_id", userID),
			zap.String("remote_addr", conn.RemoteAddr().String()))
	}()
	del.log.Debug("Websocket connected", zap.Int("user_id", userID),
		zap.String("remote_addr", conn.RemoteAddr().String()))

	return del.handleConnection(conn, userID)
}

func (del *delivery) handleConnection(conn *ws.Conn, userID int) error {
	del.connService.AddConnection(userID, conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			del.log.Debug("Read from connection error", zap.Error(err), zap.Int("user_id", userID),
				zap.String("remote_addr", conn.RemoteAddr().String()))
			del.connService.RemoveConnection(userID, conn)
			return nil
		}

		var msgReq msgRequest
		err = msgReq.UnmarshalJSON(message)
		if err != nil {
			del.log.Debug("Failed to unmarshal message", zap.Error(err), zap.Int("user_id", userID),
				zap.String("remote_addr", conn.RemoteAddr().String()))

			errResp := errorResponse{Type: "error", ErrMsg: "invalid json", ErrCode: 1}
			data, err := errResp.MarshalJSON()
			if err != nil {
				del.log.Error("Marshal json failed", zap.Error(err), zap.Int("user_id", userID),
					zap.String("remote_addr", conn.RemoteAddr().String()), zap.Any("message", errResp))

				del.connService.RemoveConnection(userID, conn)
				return nil
			}

			err = conn.WriteMessage(ws.TextMessage, data)
			if err != nil {
				del.log.Debug("Write json failed", zap.Error(err), zap.Int("user_id", userID),
					zap.String("remote_addr", conn.RemoteAddr().String()), zap.Any("message", errResp))

				del.connService.RemoveConnection(userID, conn)
				return nil
			}

			continue
		}
		del.log.Debug("Got message", zap.Int("user_id", userID),
			zap.String("remote_addr", conn.RemoteAddr().String()), zap.Int("receiver_id", msgReq.ReceiverID),
			zap.String("text", msgReq.Text))

		var sendNewChat bool
		chat, err := del.serv.GetByUsers(userID, msgReq.ReceiverID)
		if err != nil {
			if errors.Is(err, pkgErrors.ErrChatNotFound) {
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

		params := pkgChats.CreateMessageParams{AuthorID: userID, ChatID: chat.ID, Text: msgReq.Text}
		createdMessage, err := del.serv.CreateMessage(&params)
		if err != nil {
			return err
		}

		if sendNewChat {
			err = del.sendNewChatToChatMembers(chat)
			if err != nil {
				return err
			}
		}

		go func() {
			err := del.sendMessageToChatMembers(*createdMessage, chat.User1ID, chat.User2ID)
			if err != nil {
				del.log.Error(constants.FailedSendMessageToChatMembers, zap.Error(err))
			}
		}()
	}
}
