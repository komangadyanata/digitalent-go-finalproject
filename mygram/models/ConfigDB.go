package models

type ConfigMySQL struct {
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
	DBCharset string
}

type ConfigPostgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
