package entities

const (
	Ok                      = 0    // everything is ok
	ErrTooManyRequests      = 429  // too many requests
	ErrUnknown              = 1000 // something wrong
	ErrInvalidParams        = 1001 // request params invalid
	ErrInvalidAuthHeader    = 1002 // request auth failed
	ErrInvalidInitData      = 1003 // request init data invalid
	ErrParseInitDataFailed  = 1004 // parse init data failed
	ErrGenUserSessionFailed = 1005 // generate user session and set to cache failed
	ErrUserNotFound         = 1006 // not found user in database or cache
	ErrUserAuthExpired      = 1007 // user auth expired

	// ErrInternalDBInsertFailed start Internal error code
	ErrInternalDBInsertFailed      = 2000
	ErrInternalDBQueryFailed       = 2001
	ErrInternalDBUpdateFailed      = 2002
	ErrInternalDBDeleteFailed      = 2003
	ErrInternalGenerateTokenFailed = 2004

	ErrProxyCreateRequestFailed = 3001
	ErrProxyRequestFailed       = 3002
	ErrProxyParseResBodyFailed  = 3003
	ErrProxyReadResBodyFailed   = 3004
)
