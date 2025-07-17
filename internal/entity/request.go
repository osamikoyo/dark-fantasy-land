package entity

import "time"

type Request struct {
	CensoredAt  time.Time
	Payload     map[string]string
	Censored    bool
	Description string
}
