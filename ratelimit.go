package ratelimit

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/acoshift/middleware"
)

// Config is the retelimit config
type Config struct {
	Rate  int64
	Burst int
	Per   time.Duration
	Store Store
}

// New creates new ratelimit middleware from config
func New(config Config) middleware.Middleware {
	if config.Store == nil {
		panic("ratelimit: nil store")
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k := r.RemoteAddr + "|" + strconv.FormatInt(time.Now().Truncate(config.Per).UnixNano(), 10)
			c, _ := config.Store.Incr(k, config.Per)
			fmt.Println(k, c)

			if c > config.Rate {
				fmt.Println("too many")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
