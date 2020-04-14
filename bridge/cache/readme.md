## Cache

Cache is package use for caching mechanism, for this package can handle many package as driver implementior.
for example below use go redis driver.

```go

...
type SampleRedis struct {
    Abc string
}
		
...

redis := cache.New(cache.Config{
	Host:     "localhost",
	Port:     6379,
	Password: "",
}, cache.GoRedis)

if err := redis.SetNX("test_value", SampleRedis{Abc:"asa"}, 10 * time.Minute); err != nil {
	log.Println("check error", err)
}
abc := SampleRedis{}

if err := redis.Get("test_value", &abc); err != nil {
	log.Println("check error", err)
}

log.Println(abc.Abc)

```