package entity

import "time"

type JwtData struct {
	ID         int
	UserRole   string
	Expiration time.Time
}
