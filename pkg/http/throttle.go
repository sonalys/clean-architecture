package http

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
	go i.IPCleanupRoutine()
	return i
}

func (i *IPRateLimiter) IPCleanupRoutine() {
	for {
		i.mu.Lock()
		for index, ip := range i.ips {
			if ip.Reserve().Delay() > time.Minute {
				delete(i.ips, index)
			}
		}
		i.mu.Unlock()
		time.Sleep(time.Minute)
	}
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(i.r, i.b)

	i.ips[ip] = limiter

	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}
