package with

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type counter struct {
	mu     sync.RWMutex
	counts map[string]int64
	sink   chan string
}

func newCounter() *counter {
	c := counter{}
	c.counts = map[string]int64{}
	c.sink = make(chan string)
	return &c
}

func (c *counter) start() *counter {
	go func() {
		for r := range c.sink {
			c.mu.Lock()
			c.counts[r]++
			c.mu.Unlock()
		}
	}()
	return c
}

func (c *counter) string() string {
	b := strings.Builder{}
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.counts {
		fmt.Fprintf(&b, "%s\t%d\n", k, v)
	}
	return b.String()
}

func Stats(h http.Handler) http.Handler {
	c := newCounter().start()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.URL.Query()["stats"]; ok {
			w.Write([]byte(c.string()))
			return
		}
		go func() { c.sink <- r.RemoteAddr }()
		h.ServeHTTP(w, r)
	})
}
