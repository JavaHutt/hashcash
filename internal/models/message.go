package models

type messageKind string

const (
	MessageKindRequestChallenge = "request-challenge"
	MessageKindSolvedChallenge  = "solved-challenge"
)

// Message is an object that client sends to the server
type Message struct {
	Kind     messageKind `json:"kind"`
	Hashcash string      `json:"hashcash"`
}

func (k messageKind) String() string {
	return string(k)
}
