# Codephile

### Environment Variables

Environment variable is a way to store/pass some sensitive/config information that is required by the software. This can include passwords, secret keys, config variables.

To setup environment variables, create a `.env` file at conf directory of project containing following information:
```
PORT = <The port to be used>
DBPath = <Connection string of local database>
DBName = <Name of database>
HMACKEY = <HMAC Encryption key>
REDISURL = <URL of redis server>
REDISPASSWD = <Redis Password(empty if no password)>
```
