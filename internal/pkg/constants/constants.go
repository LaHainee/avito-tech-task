package constants

import "time"

const (
	ConfigPath = "config/config.toml"

	InvalidBodyMessage   = "Invalid body"
	InvalidUserIDMessage = "Invalid user id"
	InvalidQueryParams   = "Invalid query params"

	// CurrencyAPIUpdatePeriod = 24 * time.Hour
	CurrencyAPIUpdatePeriod = 10 * time.Second
)
