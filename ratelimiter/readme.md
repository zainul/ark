## Rate Limiter

How to use
```go

// create the rules
rules := []Rule{}

rules = append(
    rules, 
    rule.NewSlidingWindow( 
        "some/api", 
        60, 
        rule.SlidingWindowLimiterByRPS, 
        1, 
        rds,
    )
)

// crate rate limiter
rule := New(rules)

// use in rate limit
isRateLimit, err := rule.IsRateLimit(time.Now())

```