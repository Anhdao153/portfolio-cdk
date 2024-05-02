package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	CDK_DEFAULT_ACCOUNT string
	CDK_DEFAULT_REGION  string
	SOURCE_BE           string
	MaxDuration			int
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	envFilePath := filepath.Join(cwd, "/.env")
	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("ENV got something error: %v", err)
	}
	CDK_DEFAULT_ACCOUNT = os.Getenv("CDK_DEFAULT_ACCOUNT")
	CDK_DEFAULT_REGION = os.Getenv("CDK_DEFAULT_REGION")
	SOURCE_BE = os.Getenv("SOURCE_BE")
	MaxDuration = 60 // second
}
