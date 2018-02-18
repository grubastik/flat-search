# flat-search

This is a small project for look-up a flat according to the options stored in config file in CZ.

This cli should be installed to run migrations:
```
https://github.com/mattes/migrate
```
Note: to run migrations please use the following command:
```bash
$GOPATH/bin/migrate -url mysql://[DB_USER_NAME]:[DB_USER_PASSWORD]@tcp\([IP_OR_HOSTNAME]:[PORT]\)/[DB_NAME] -path ./migrations/ up
```
Note2: before using migrations it is required to create DB and grant all necessary permissions to your DB user.

## TODO 
1. Minimize usage of pointers
2. tests
3. move sreality to separate repo
