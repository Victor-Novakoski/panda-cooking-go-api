package ratelimit_test

import (
	"testing"
	"time"

	"panda-cooking-go-api/internal/ratelimit"

	"github.com/stretchr/testify/assert"
)

// newLimiter cria um limitador com relógio controlado pelo teste.
func newLimiter(max int, window time.Duration) (*ratelimit.Limiter, *time.Time) {
	now := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	l := ratelimit.New(max, window)
	l.Now = func() time.Time { return now }
	return l, &now
}

func TestLimiter_Allow(t *testing.T) {
	t.Run("aceita até o limite e recusa a seguinte", func(t *testing.T) {
		l, _ := newLimiter(3, time.Minute)

		for range 3 {
			ok, _ := l.Allow("1.2.3.4")
			assert.True(t, ok)
		}
		ok, wait := l.Allow("1.2.3.4")

		assert.False(t, ok)
		assert.Equal(t, time.Minute, wait)
	})

	t.Run("chaves diferentes têm contagens separadas", func(t *testing.T) {
		l, _ := newLimiter(1, time.Minute)

		ok1, _ := l.Allow("a")
		ok2, _ := l.Allow("b")

		assert.True(t, ok1)
		assert.True(t, ok2)
	})

	t.Run("libera de novo quando a janela acaba", func(t *testing.T) {
		l, now := newLimiter(1, time.Minute)
		l.Allow("a")
		ok, _ := l.Allow("a")
		assert.False(t, ok)

		*now = now.Add(time.Minute)
		ok, _ = l.Allow("a")

		assert.True(t, ok)
	})
}

func TestLimiter_BlockedAndReset(t *testing.T) {
	t.Run("bloqueia ao chegar no limite sem contar tentativa", func(t *testing.T) {
		l, now := newLimiter(2, 15*time.Minute)

		blocked, _ := l.Blocked("v@v.com")
		assert.False(t, blocked)

		l.Allow("v@v.com")
		l.Allow("v@v.com")
		*now = now.Add(5 * time.Minute)
		blocked, wait := l.Blocked("v@v.com")

		assert.True(t, blocked)
		assert.Equal(t, 10*time.Minute, wait)
	})

	t.Run("reset zera a contagem", func(t *testing.T) {
		l, _ := newLimiter(1, time.Minute)
		l.Allow("a")

		l.Reset("a")
		blocked, _ := l.Blocked("a")

		assert.False(t, blocked)
	})
}
