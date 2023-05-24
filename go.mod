module github.com/mdg-iitr/Codephile

// +heroku install ./cmd/... .

go 1.12

require (
	cloud.google.com/go/storage v1.10.0
	firebase.google.com/go v3.9.0+incompatible
	github.com/astaxie/beego v1.12.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/sentry-go v0.5.1
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gocolly/colly v1.2.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/schema v1.1.0
	github.com/joho/godotenv v1.3.0
	github.com/smartystreets/goconvey v1.6.4
	go.mongodb.org/mongo-driver v1.11.6
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/oauth2 v0.0.0-20210413134643-5e61552d6c78
	google.golang.org/api v0.30.0
)

require (
	cloud.google.com/go/firestore v1.1.0 // indirect
	github.com/PuerkitoBio/goquery v1.5.0 // indirect
	github.com/antchfx/htmlquery v1.1.0 // indirect
	github.com/antchfx/xmlquery v1.1.0 // indirect
	github.com/antchfx/xpath v1.1.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/temoto/robotstxt v1.1.1 // indirect
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/tools v0.0.0-20201211185031-d93e913c1a58 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
