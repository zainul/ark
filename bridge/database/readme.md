## Database

Database is package use for to connect and interaction with any dbms with many driver implementor.
For example in this package we can inject package using `GORM` or `SQLT`

```go

...

masterDB, slaveDB := initial.InitDBMS(cfg)

...

// INIT DB
db := database.NewDB("postgres", database.Config{
    MasterConn: masterDB,
	SlaveConn:  slaveDB,
}, database.GORM)

...

result := make([]entity.User, 0)
err := s.db.EntityBy(field, value, &result)

```