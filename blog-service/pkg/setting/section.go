package setting

import "time"

type ServerSettingS struct {
	RunMode string
	HttpPort string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize int
	MaxPageSize int

	LogSavePath string
	LogFileName string
	LogFileExt string

	UploadSavePath string
	UploadServerUrl string
	UploadImageMaxSize int
	UploadImageAllowExts []string
}

type DatabaseSettingS struct {
	DBType string
	Username string
	Password string
	Host string
	Port int
	DBName string
	TablePrefix string
	Charset string
	Parsetime bool
	MaxIdleConns int
	MaxOpenConns int
}

type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

type EmailSettingS struct {
	Host string
	Port int
	UserName string
	Password string
	IsSSL bool
	From string
	To []string
}