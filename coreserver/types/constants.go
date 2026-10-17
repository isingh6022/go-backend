package types

import (
	"time"
)

const DefaultPanicMessage = "internal server error"

const DefaultShutdownTimeout = 15 * time.Second

const DefaultReadHeaderTimeout = 5 * time.Second
const DefaultReadTimeout = 15 * time.Second
const DefaultWriteTimeout = 30 * time.Second
const DefaultIdleTimeout = 120 * time.Second
