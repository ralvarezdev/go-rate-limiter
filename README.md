# go-rate-limiter

Redis-backed rate limiter for Go. It limits the number of requests per key (for example a client IP) within a time period. Requires Go 1.24 (per `go.mod`).

## Installation

```bash
go get github.com/ralvarezdev/go-rate-limiter
```

Dependencies include `go-redis/v9`, `go-databases` and `go-strings`.

## API (package `redis`)

- **`RateLimiter`** — interface with `Limit(ip string) error`.
- **`NewDefaultRateLimiter(redisClient, limit, period)`** — returns an error when the client is nil. The first request creates the counter with `period` as TTL; later requests increment it and return `ErrTooManyRequests` once the count reaches `limit`.
- **`GetKey(ip)`** — builds `rate_limiter:<ip>`; `KeyPrefix` and `KeySeparator` are exported variables. `SetInitialValue(key)` initialises the counter.
- **Errors** — `ErrTooManyRequests`, `ErrNilRateLimiter`.

## Usage

```go
client := redis.NewClient(&redis.Options{Addr: redisAddr})

limiter, err := goratelimiterredis.NewDefaultRateLimiter(client, 100, time.Minute)
if err != nil {
    return err
}

if err = limiter.Limit(clientIP); errors.Is(err, goratelimiterredis.ErrTooManyRequests) {
    // respond with HTTP 429
}
```

The `go-gin` and `go-net` repositories provide HTTP middleware that consumes this limiter.

## Development

```bash
go build ./...
```

There are no tests.

## License

GNU General Public License v3.0 (see [LICENSE](LICENSE)).
