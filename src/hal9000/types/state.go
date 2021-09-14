package types

type State interface {
	Name() string
	ProcessIncomingMessage(c Person, m RequestMessage) (State, ResponseMessage, error)
}
