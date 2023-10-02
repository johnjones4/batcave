package push

import (
	"context"
	"main/core"
)

type SocketSender struct {
	activeSockets map[string]map[string]*core.ActiveSocket
}

func NewSocketSender() *SocketSender {
	return &SocketSender{
		activeSockets: make(map[string]map[string]*core.ActiveSocket),
	}
}

func (s *SocketSender) RegisterActiveSocket(clientId string, connectionId string) *core.ActiveSocket {
	ok, _ := s.DeregisterActiveSocket(clientId, connectionId)
	if !ok {
		s.activeSockets[clientId] = make(map[string]*core.ActiveSocket)
	}
	as := &core.ActiveSocket{
		Messages:     make(chan core.PushMessage, 255),
		ClientId:     clientId,
		ConnectionId: connectionId,
	}
	s.activeSockets[clientId][connectionId] = as
	return as
}

func (s *SocketSender) DeregisterActiveSocket(clientId string, connectionId string) (bool, bool) {
	var ok2 bool
	_, ok1 := s.activeSockets[clientId]
	if ok1 {
		_, ok2 = s.activeSockets[clientId][connectionId]
		if ok2 {
			delete(s.activeSockets[clientId], connectionId)
		}
	}
	return ok1, ok2
}

func (s *SocketSender) SendToClient(ctx context.Context, clientId string, message core.PushMessage) (bool, error) {
	sockets, ok := s.activeSockets[clientId]
	if !ok {
		return false, nil
	}
	for _, socket := range sockets {
		socket.Messages <- message
	}
	return len(sockets) > 0, nil
}
