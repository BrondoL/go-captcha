package constant

import "time"

const (
	CacheKeyCaptcha = "CAPTCHA"
)

const (
	CacheTTLOneDay     = 24 * time.Hour
	CacheTTLOneMinute  = 1 * time.Minute
	CacheTTLForever    = 0
	CacheTTLInvalidate = -1
)
