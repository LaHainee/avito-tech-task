package constants

import "time"

const (
	_ = iota
	ADD
	REDUCE
	TRANSFER

	ConfigPath              = "config/config.toml"
	InvalidBodyMessage      = "Invalid body"
	InvalidUserIDMessage    = "Invalid user id"
	InvalidQueryParams      = "Invalid query params"
	CurrencyAPIUpdatePeriod = 24 * time.Hour
)
