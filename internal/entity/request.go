package entity

import "time"

type Request[T any] struct {
	CensoredAt  time.Time
	Payload     T
	Status      string
	Description string
}
