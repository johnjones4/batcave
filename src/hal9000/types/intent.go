package types

type Intent interface {
	Execute(runtime Runtime, lastState State) (State, ResponseMessage, error)
}
