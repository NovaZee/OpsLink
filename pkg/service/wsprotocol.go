package service

import (
	"time"
)

const (
	pingFrequency = 10 * time.Second
	pingTimeout   = 2 * time.Second
)

func WSSignalConnection() {
}
