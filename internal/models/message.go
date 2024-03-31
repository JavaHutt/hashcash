package models

type messageKind string

const (
	MessageKindRequestChallenge = "request-challenge"
	MessageKindSolvedChallenge  = "solved-challenge"
)

type Message struct {
	Kind     messageKind `json:"kind"`
	Hashcash string      `json:"hashcash"`
}
