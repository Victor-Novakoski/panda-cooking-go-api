// Package ratelimit conta tentativas por chave (IP, e-mail) em janela fixa,
// em memória. Serve para uma instância só da API; com várias, o contador
// precisa ir para um armazenamento compartilhado.
package ratelimit

import (
	"sync"
	"time"
)

type entry struct {
	count int
	reset time.Time
}

type Limiter struct {
	mu        sync.Mutex
	max       int
	window    time.Duration
	hits      map[string]*entry
	lastSweep time.Time

	// Now pode ser trocado nos testes para controlar o relógio.
	Now func() time.Time
}

// New cria um limitador que aceita até max tentativas por chave a cada window.
func New(max int, window time.Duration) *Limiter {
	return &Limiter{max: max, window: window, hits: map[string]*entry{}, Now: time.Now}
}

// Allow conta uma tentativa para a chave. Devolve false e quanto falta para
// liberar quando a chave já passou do limite na janela atual.
func (l *Limiter) Allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.Now()
	l.sweep(now)

	e := l.current(key, now)
	if e == nil {
		e = &entry{reset: now.Add(l.window)}
		l.hits[key] = e
	}
	e.count++
	if e.count > l.max {
		return false, e.reset.Sub(now)
	}
	return true, 0
}

// Blocked diz se a chave já chegou ao limite, sem contar nova tentativa.
func (l *Limiter) Blocked(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.Now()
	if e := l.current(key, now); e != nil && e.count >= l.max {
		return true, e.reset.Sub(now)
	}
	return false, 0
}

// Reset zera a contagem da chave (ex.: login certo).
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.hits, key)
}

// current devolve a entrada da chave se a janela dela ainda vale.
func (l *Limiter) current(key string, now time.Time) *entry {
	e, ok := l.hits[key]
	if !ok || !now.Before(e.reset) {
		return nil
	}
	return e
}

// sweep apaga as janelas vencidas, no máximo uma vez por janela,
// para o mapa não crescer sem limite.
func (l *Limiter) sweep(now time.Time) {
	if now.Sub(l.lastSweep) < l.window {
		return
	}
	for k, e := range l.hits {
		if !now.Before(e.reset) {
			delete(l.hits, k)
		}
	}
	l.lastSweep = now
}
