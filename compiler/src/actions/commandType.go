package actions

type action = byte

type command struct {
	Action action
}

func isAction() bool {
	return false
}
