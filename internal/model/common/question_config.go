package common

import "time"

type Config struct {
	TimeLimit   time.Duration `json:"time_limit,omitempty"`
	MemoryLimit uint64        `json:"memory_limit,omitempty"`
}
