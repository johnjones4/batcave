package types

type State interface {
	Name() string
	ProcessIncomingMessage(r Runtime, c Person, m RequestMessage) (State, ResponseMessage, error)
}
