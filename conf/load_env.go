package conf

import (
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var AppRootDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	AppRootDir = filepath.Dir(filepath.Dir(file))
	err := godotenv.Load(filepath.Join(filepath.Dir(file), ".env"))
	if err != nil {
		log.Println("No .env file found")
	}
	search.GetElasticClient()
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}
