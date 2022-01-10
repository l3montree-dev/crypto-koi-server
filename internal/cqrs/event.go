package cqrs

type Event interface {
	GetId() string
	GetType() string
	ToJSON() string
}
