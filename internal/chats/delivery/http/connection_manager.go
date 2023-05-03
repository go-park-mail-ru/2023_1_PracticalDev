package http

import (
	"fmt"
	"sync"

	ws "github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
)

type ConnectionManager struct {
	connections map[int][]*ws.Conn // UserID -> []Conn
	mu          sync.Mutex
	log         log.Logger
}

func NewConnectionManager(log log.Logger) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[int][]*ws.Conn),
		log:         log,
	}
}

func (cm *ConnectionManager) AddConnection(userID int, conn *ws.Conn) {
	cm.mu.Lock()
	cm.connections[userID] = append(cm.connections[userID], conn)
	cm.mu.Unlock()
}

func (cm *ConnectionManager) RemoveConnection(userID int, conn *ws.Conn) {
	cm.mu.Lock()
	conns := cm.connections[userID]
	var idx int
	for idx = range conns { // find conn index to remove
		if conns[idx] == conn {
			break
		}
	}
	conns[idx] = conns[len(conns)-1]
	cm.connections[userID] = conns[:len(conns)-1]
	cm.mu.Unlock()
	cm.log.Debug(fmt.Sprintf("Connection %p was deleted from connections list", conn))
}

// Broadcast отправляет сообщение всем соединениям пользователя
func (cm *ConnectionManager) Broadcast(response any, userID int) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var err error
	conns, ok := cm.connections[userID]
	if ok {
		for _, conn := range conns {
			cm.log.Debug(fmt.Sprintf("New response %v for sending to userID=%d, conn=%p:", response, userID, conn))
			err = conn.WriteJSON(response)
			if err != nil {
				cm.log.Error("Broadcast: error: %v", err)
				return err
			}
		}
	} else {
		cm.log.Debug("There are no connections for userID=", userID)
	}

	return nil
}
