package env

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strings"
)

// Environment constants
const (
	Development, DevelopmentEnv = "development", "development"
	Staging, StagingEnv         = "staging", "staging"
	Production, ProductionEnv   = "production", "production"
)

//GetEnvMode is represent of environment mode currently applied
func GetEnvMode() string {
	return os.Getenv("ENV_MODE")
}

//GetAppPath is represent of get the path of application service
func GetAppPath() string {
	return os.Getenv("APP_PATH")
}

// Env related var
var (
	envName   = "TKPENV"
	goVersion string
)

func init() {
	// env package will read .env file when applicatino is started
	err := SetFromEnvFile(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("failed to set env file: %v\n", err)
	}
	goVersion = runtime.Version()
}

// SetFromEnvFile read env file and set the environment variables
func SetFromEnvFile(filepath string) error {
	if _, err := os.Stat(filepath); err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		return err
	}
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		vars := strings.SplitN(text, "=", 2)
		if len(vars) < 2 {
			return err
		}
		if err := os.Setenv(vars[0], vars[1]); err != nil {
			return err
		}
	}
	return nil
}

// ServiceEnv return TKPENV service environment
func ServiceEnv() string {
	e := os.Getenv(envName)
	if e == "" {
		e = DevelopmentEnv
	}
	return e
}

// GoVersion to return current build go version
func GoVersion() string {
	return goVersion
}

// IsDevelopment return true when env is "development"
func IsDevelopment() bool {
	return ServiceEnv() == DevelopmentEnv
}

// IsStaging return true when env is "staging"
func IsStaging() bool {
	return ServiceEnv() == StagingEnv
}

// IsProduction return true when env is "production"
func IsProduction() bool {
	return ServiceEnv() == ProductionEnv
}
