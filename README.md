# Codephile | [![CircleCI](https://circleci.com/gh/mdg-iitr/Codephile.svg?style=svg&circle-token=f989c04ad5d3a6578d45296b18cdca223e504bde)](https://circleci.com/gh/mdg-iitr/Codephile)
## Services
We use the following services in our server,

* MongoDB: Main database of the server, stores user info, submission, profile,etc. Install from [here](https://docs.mongodb.com/manual/installation/)
* Redis: Used to logout and blacklist users. Serves as cache for contests API. Download from [here](https://redis.io/download) 
* Elastic Search: Some of user data is indexed to elasticsearch db, in order to use the search API. Download it from [here](elastic.co/downloads/)
* Firebase storage: The profile pictures are stored in firebase storage. Create a firebase account.

## Environment Variables

Environment variable is a way to store/pass some sensitive/config information that is required by the software. This can include passwords, secret keys, config variables.

To setup environment variables, create a `.env` file at conf directory of project containing following information:
```
PORT=<The port to be used: optional>
DBPath=<Connection string of local database>
HMACKEY=<HMAC Encryption key>
REDISURL=<connection string of redis server>
FIREBASE_CONFIG=<Firebase config including bucket name(json)>
FIREBASE_CREDENTIALS=<Firebase admin SDK credentials(json)>
ELASTICURL=<connection string of elasticsearch cloud>
SENTRY_DSN=<Data source name of sentry server: optional>
EMAIL_SMTP_HOST=<host of smtp server>
EMAIL_SMTP_PORT=<port of smtp server>
EMAIL_SERVER_USER=<username of email account>
EMAIL_SERVER_PASS=<password of email account>
```
NOTE: Before proceeding further, ensure that your local .env file is present with above configuration variables.

Ask for codechef creds from the maintainer
```
CLIENT_ID=<codechef id>
CLIENT_SECRET=<codechef secret>
```

## Setup Instructions

Download golang from [here](https://golang.org/dl/) and setup GOPATH

In order to ensure the GOPATH environment variable is setup run: 
```shell script
$ echo $GOPATH
```
This should give non empty output

Now clone the repo in the appropriate directory.
```shell script
$ mkdir -p $GOPATH/src/github.com/mdg-iitr/Codephile && cd $_ 
$ git clone https://github.com/mdg-iitr/Codephile.git
```
We used beego framework to bootstrap the project. Download and setup bee command line program from [here](https://beego.me/quickstart).

In order to generate documentation from comments, run:
```shell script
$ bee run -downdoc=true -gendoc=true
```
If you didn't make any changes in documentation comment, simply run:
```shell script
$ bee run
```
Custom programs could be run using
```shell script
$ go run cmd/<path to main package go file>
```
E.g.
```shell script
 $ go run cmd/blacklist-user/blacklist_user.go
```
## Setup using docker
You can use the `dev_docker-compose.yml` file to spin up containers with Mongo, Redis & Elastic Search services easily.
Use these env variables
```
REDISURL=redis://redis:6379
ELASTICURL=http://elastic:secret@elasticsearch:9200/codephile/?sniff=false
DBPath=mongodb://mongoadmin:secret@mongo:27017/admin
```
And run these commands
```shell script
$ mkdir -m 777 -p data/elasticsearch
$ docker-compose -f dev_docker-compose.yml up
```

## Tests

Change the `DBPath` and `ELASTICURL` in .env file 

Run the tests
```shell script
$ go test -v ./tests
```

## Components

* `cmd`: Contains standalone programs for specific tasks like updating user submissions, deleting, blacklist users.

* `conf`: Contains global app level constants and configuration files. This package has to be imported first in the main package, as it loads various global variables and inits various clients(sentry, elasticsearch).

* `controller`:  Responsible for handling the requests corresponding to various endpoints. Contains separate files for separate namespaces.

* `errors`: Contains custom error messages and json response structs to respond with, in case of errors.

* `middleware`: Sits before controllers. Mainly authenticates user and extracts uid from user token

* `models`:
    * `models/db`: Handles db connection and manages connection pool. Provides a clean interface to establish db connections.
    * `models/types`: Contains the types for various database schema and response models.
    * `/`: Contains database operations, queries.
    
* `routers`: Registers endpoints. Beego generates the routes from comments inside controllers. See [this](https://beego.me/docs/mvc/controller/router.md#annotations) for more information.

* `scrappers`: Contains the main logic for scrapping user data(submission, profile) from platforms. Each platform's logic is contained in packages with the platform name and a simple interface to scrappers is exposed through `interface.go` 

* `services`: Creates and exposes the clients for various services like redis, elasticsearch. Also contains code for worker routines that are activated on request to POST `/user/submission`

* `swagger`: Contains the static files and `swagger.json` and `swagger.yml` for API documentation. Documentation could be generated using bee command line tool `bee run -downdoc=true -gendoc=true`

* `test`: Will contain tests for various endpoints and unit tests. Currently, only test for `/user/all` is present. Run the tests using 
`go test  ./tests/...` 

Beginners are advised to begin with writing some tests.

## CI

When a pull request is submitted, continuous integration jobs are run automatically to ensure the code builds and is relatively well-written. The jobs are run on circleci.
At present, the build, tests and linters are run on CI.

We use [golang-ci](https://github.com/golangci/golangci-lint) lint for linting jobs. Download and run the linter locally before submitting a PR.