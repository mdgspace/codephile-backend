package conf

import (
	"github.com/joho/godotenv"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"log"
	"path/filepath"
	"runtime"
)

var AppRootDir string

func init() {
	search.GetElasticClient()
	_, file, _, _ := runtime.Caller(0)
	AppRootDir = filepath.Dir(filepath.Dir(file))
	err := godotenv.Load(filepath.Join(filepath.Dir(file), ".env"))
	if err != nil {
		log.Println("No .env file found")
	}
}
