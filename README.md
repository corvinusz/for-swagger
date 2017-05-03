# For Swagger
Special repo for testing swagger docs generation

https://github.com/labstack/echo as web server and

https://github.com/go-xorm/xorm as a database ORM

# Installation
## Dependencies
There is still no vendoring here. Waiting for https://github.com/golang/dep

So ...

```bash
go get -u github.com/BurntSushi/toml
go get -u github.com/labstack/echo
go get -u github.com/labstack/echo/middleware
go get -u github.com/mattn/go-sqlite3
go get -u github.com/go-xorm/xorm
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/onsi/ginkgo
go get -u github.com/onsi/gomega
go get -u github.com/go-resty/resty
go get -u github.com/go-testfixtures/testfixtures
```

Hope that nothing is missed.

## Application
```bash
go install github.com/corvinusz/for-swagger
```

## Database
Currently is using *sqlite3*-database, located in file /tmp/for-swagger.sqlite.db

#License
MIT
