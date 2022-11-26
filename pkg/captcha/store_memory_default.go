package captcha

import "time"

var (
	GCLimitNumber   = 10240
	Expiration      = 10 * time.Minute
	DefaultMemStore = NewMemoryStore(GCLimitNumber, Expiration)
)
