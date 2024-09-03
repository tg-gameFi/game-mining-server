package configs

type BasicConfig struct {
	Env               string `json:"env"`               // running env: dev,test,prod
	Port              int    `json:"port"`              // server port: 8080
	Version           int    `json:"version"`           // api version: 1
	SessionExpiresSec int    `json:"sessionExpiresSec"` // session expires time unit in seconds
	SessionEncryptKey string `json:"sessionEncryptKey"` // session encrypt key
	CheckinBrokenSec  int    `json:"checkinBrokenSec"`  // checkin continuous broken duration time in seconds
}

type DatabaseConfig struct {
	Host     string `json:"host"`     // db hostname
	Port     string `json:"port"`     // db port
	User     string `json:"user"`     // db username
	Pass     string `json:"pass"`     // db password
	Dbname   string `json:"dbname"`   // db name
	InitPath string `json:"initPath"` // init.sql file path, if not null, will exec init.sql to create dbs and tables
}

// CacheConfig Cache config, use Redis by default
type CacheConfig struct {
	Host string `json:"host"` // cache server hostname
	Port string `json:"port"` // cache server port
	Pass string `json:"pass"` // cache server password, set to empty means need tls
}

type BotConfig struct {
	WebUrl string `json:"webUrl"`
}

type Config struct {
	Basic    *BasicConfig    `json:"basic"`
	Database *DatabaseConfig `json:"database"`
	Cache    *CacheConfig    `json:"cache"`
	Bot      *BotConfig      `json:"bot"`
}
