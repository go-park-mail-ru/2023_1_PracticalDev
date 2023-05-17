package connectionservice

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	"sync"

	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Service struct {
	connections map[int][]*ws.Conn // UserID -> Connections
	mu          sync.Mutex
	log         *zap.Logger
}

func NewService(log *zap.Logger) *Service {
	return &Service{
		connections: make(map[int][]*ws.Conn),
		log:         log,
	}
}

func (serv *Service) AddConnection(userID int, conn *ws.Conn) {
	serv.mu.Lock()
	serv.connections[userID] = append(serv.connections[userID], conn)
	serv.mu.Unlock()
	serv.log.Debug("Connection was added", zap.Int("user_id", userID),
		zap.String("remote_addr", conn.RemoteAddr().String()))
}

func (serv *Service) RemoveConnection(userID int, conn *ws.Conn) {
	serv.mu.Lock()
	conns := serv.connections[userID]
	var idx int
	for idx = range conns { // find conn index to remove
		if conns[idx] == conn {
			break
		}
	}
	conns[idx] = conns[len(conns)-1]
	serv.connections[userID] = conns[:len(conns)-1]
	serv.mu.Unlock()
	serv.log.Debug("Connection was removed", zap.Int("user_id", userID),
		zap.String("remote_addr", conn.RemoteAddr().String()))
}

// Broadcast отправляет сообщение всем соединениям пользователя
func (serv *Service) Broadcast(response any, userID int) error {
	serv.mu.Lock()
	defer serv.mu.Unlock()

	var err error
	conns, ok := serv.connections[userID]
	if ok {
		for _, conn := range conns {
			err = conn.WriteJSON(response)
			if err != nil {
				serv.log.Error(constants.FailedWriteJSONResponse, zap.Error(err))
				return err
			}
			serv.log.Debug("Response was sent", zap.Int("user_id", userID),
				zap.String("remote_addr", conn.RemoteAddr().String()), zap.Any("response", response))
		}
	}

	return nil
}
