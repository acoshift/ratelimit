package ratelimit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/acoshift/ratelimit"
)

func repeat(h http.Handler, times int) {
	for i := 0; i < times; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		h.ServeHTTP(w, r)
	}
}

func TestRate(t *testing.T) {
	var called int

	h := ratelimit.New(ratelimit.Config{
		Rate:  10,
		Per:   time.Second,
		Store: ratelimit.NewMemoryStore(),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
	}))

	t.Run("MustNotLimit", func(t *testing.T) {
		called = 0
		repeat(h, 10)
		time.Sleep(time.Second)
		repeat(h, 10)
		if called != 20 {
			t.Fatalf("expected allow 20 called; got %d", called)
		}
	})

	time.Sleep(time.Second)

	t.Run("MustLimit", func(t *testing.T) {
		called = 0
		repeat(h, 20)
		time.Sleep(time.Second)
		repeat(h, 20)
		if called != 20 {
			t.Fatalf("expected allow 20 called; got %d", called)
		}
	})
}
