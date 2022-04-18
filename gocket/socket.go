package gocket

type socket struct {
	IRoom
}

func Socket() *socket {
	return &socket{}
}
