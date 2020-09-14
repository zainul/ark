package env

// Environment constants
const (
	Development = "development"
	Staging     = "staging"
	Production  = "production"
)

//GetEnvMode is represent of environment mode currently applied
func GetEnvMode() string {
	return os.Getenv("ENV_MODE")
}

//GetAppPath is represent of get the path of application service
func GetAppPath() string {
	return os.Getenv("APP_PATH")
}