package server

import (
	"hash/fnv"
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	requestsPerWindow = 240
	rateWindow        = time.Minute
	maxClientsPerShard = 1024
)

type rateLimitShard struct {
	mu      sync.Mutex
	clients map[string]rateLimitClient
	checks  uint16
}

type rateLimitClient struct {
	windowStart time.Time
	requests    int
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := rateLimitClientAddress(r.RemoteAddr)
		hash := fnv.New32a()
		_, _ = hash.Write([]byte(client))
		shard := &s.limiters[hash.Sum32()%uint32(len(s.limiters))]
		now := time.Now()

		shard.mu.Lock()
		if shard.clients == nil {
			shard.clients = make(map[string]rateLimitClient)
		}
		entry, exists := shard.clients[client]
		if !exists && len(shard.clients) >= maxClientsPerShard {
			shard.checks++
			if shard.checks%128 == 0 {
				for key, candidate := range shard.clients {
					if now.Sub(candidate.windowStart) >= rateWindow {
						delete(shard.clients, key)
					}
				}
			}
			if len(shard.clients) >= maxClientsPerShard {
				shard.mu.Unlock()
				writeRateLimitError(w, int(rateWindow.Seconds()))
				return
			}
		}
		if !exists || now.Sub(entry.windowStart) >= rateWindow {
			entry = rateLimitClient{windowStart: now}
		}
		entry.requests++
		shard.clients[client] = entry
		remaining := rateWindow - now.Sub(entry.windowStart)
		limited := entry.requests > requestsPerWindow
		shard.mu.Unlock()

		if limited {
			retryAfter := int(math.Ceil(remaining.Seconds()))
			if retryAfter < 1 {
				retryAfter = 1
			}
			writeRateLimitError(w, retryAfter)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func rateLimitClientAddress(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && host != "" {
		return host
	}
	return remoteAddr
}

func writeRateLimitError(w http.ResponseWriter, retryAfter int) {
	w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	respondError(w, http.StatusTooManyRequests, "Too many requests")
}
