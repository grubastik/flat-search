# flat-search

This is a small project for look-up a flat according to the options stored in config file in CZ.

To run the project please, use the following instructions:
1. Configure mysql database server
2. Create user and password for database
3. Clone this repo using following commands:
```bash
mkdir -p ~/go/src/github.com/grubastik
cd ~/go/src/github.com/grubastik
git clone https://github.com/grubastik/flat-search.git
```
4. run `make pre-install`
5. run `make migrate`
6. run `make install`

## For docker containers
7. run `make run.docker`
8. run `make docker.start`
9. run `make docker.stop`

This cli should be installed to run migrations
```
https://github.com/mattes/migrate
```
To install migrate use:
```
go get -u github.com/golang-migrate/migrate/cli
```
Note: to run migrations please use the following command:
```bash
$GOPATH/bin/cli -url mysql://[DB_USER_NAME]:[DB_USER_PASSWORD]@tcp\([IP_OR_HOSTNAME]:[PORT]\)/[DB_NAME] -path ./migrations/ up
```
Note2: before using migrations it is required to create DB and grant all necessary permissions to your DB user.

## TODO 
1. Minimize usage of pointers
2. tests
3. move sreality to separate repo
