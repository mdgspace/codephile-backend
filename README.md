# Codephile

### Environment Variables

Environment variable is a way to store/pass some sensitive/config information that is required by the software. This can include passwords, secret keys, config variables.

To setup environment variables, create a `.env` file at conf directory of project containing following information:
```
PORT = <The port to be used>
DBPath = <Connection string of local database>
HMACKEY = <HMAC Encryption key>
REDISURL = <connection string of redis server>
FIREBASE_CONFIG = <Firebase config including bucket name(json)>
FIREBASE_CREDENTIALS = <Firebase admin SDK credentials(json)>
ELASTICURL = <connection string of elasticsearch cloud>
```
