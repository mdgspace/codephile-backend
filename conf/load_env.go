package conf

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

var AppRootDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	AppRootDir = filepath.Dir(filepath.Dir(file))
	err := godotenv.Load(filepath.Join(filepath.Dir(file), ".env"))
	if err != nil {
		log.Println(err.Error())
	}
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}
