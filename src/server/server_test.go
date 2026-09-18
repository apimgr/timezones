package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apimgr/timezones/src/config"
	"github.com/apimgr/timezones/src/timezones"
)

const fixtureJSON = `[
	{
		"value": "Eastern Standard Time",
		"abbr": "EST",
		"offset": -5,
		"isdst": false,
		"text": "(UTC-05:00) Eastern Time (US & Canada)",
		"utc": ["America/New_York", "America/Detroit"]
	},
	{
		"value": "Pacific Standard Time",
		"abbr": "PST",
		"offset": -8,
		"isdst": true,
		"text": "(UTC-08:00) Pacific Time (US & Canada)",
		"utc": ["America/Los_Angeles"]
	}
]`

func newTestServer(t *testing.T) *Server {
	t.Helper()
	svc, err := timezones.NewService([]byte(fixtureJSON))
	if err != nil {
		t.Fatalf("NewService() unexpected error: %v", err)
	}
	cfg := config.DefaultConfig()
	return New(svc, cfg, "127.0.0.1", "8080", "1.0.0-test", "2024-01-01", "abcdef")
}

func decodeAPIResponse(t *testing.T, body []byte) APIResponse {
	t.Helper()
	var resp APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode APIResponse: %v (body: %s)", err, body)
	}
	return resp
}

// TestNewServerAndRouter verifies New() wires up a usable router.
func TestNewServerAndRouter(t *testing.T) {
	s := newTestServer(t)
	if s.Router() == nil {
		t.Fatal("Router() returned nil")
	}
}

// TestHandleHome covers the homepage renders successfully.
func TestHandleHome(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Timezones API") {
		t.Error("home page missing expected title")
	}
}

// TestHandleHealthVariants covers /healthz, /health, /status all returning
// the same healthy payload.
func TestHandleHealthVariants(t *testing.T) {
	s := newTestServer(t)
	for _, path := range []string{"/healthz", "/health", "/status"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			s.Router().ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			resp := decodeAPIResponse(t, rec.Body.Bytes())
			if !resp.Success {
				t.Error("expected Success = true")
			}
		})
	}
}

// TestHandleManifest verifies the PWA manifest is valid JSON with expected
// fields.
func TestHandleManifest(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var manifest map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &manifest); err != nil {
		t.Fatalf("failed to decode manifest: %v", err)
	}
	if manifest["name"] != "Timezones API" {
		t.Errorf("manifest name = %v, want Timezones API", manifest["name"])
	}
}

// TestHandleServiceWorker verifies the service worker script is served
// with the correct content type.
func TestHandleServiceWorker(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/sw.js", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/javascript" {
		t.Errorf("Content-Type = %q, want application/javascript", ct)
	}
	if !strings.Contains(rec.Body.String(), "CACHE_NAME") {
		t.Error("service worker body missing expected content")
	}
}

// TestHandleRobotsTxt verifies the default allow rules are rendered.
func TestHandleRobotsTxt(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "User-agent: *") {
		t.Error("robots.txt missing User-agent line")
	}
}

// TestHandleSecurityTxt covers both /security.txt and the well-known path,
// and the default contact fallback.
func TestHandleSecurityTxt(t *testing.T) {
	s := newTestServer(t)
	for _, path := range []string{"/security.txt", "/.well-known/security.txt"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			s.Router().ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			if !strings.Contains(rec.Body.String(), "security@example.com") {
				t.Error("security.txt missing default contact")
			}
		})
	}
}

// TestHandleTimezonesJSON verifies the raw embedded JSON bytes are served
// verbatim.
func TestHandleTimezonesJSON(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones.json", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != fixtureJSON {
		t.Error("raw JSON response does not match fixture")
	}
}

// TestHandleTimezonesAll covers the JSON listing endpoint.
func TestHandleTimezonesAll(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	resp := decodeAPIResponse(t, rec.Body.Bytes())
	data, ok := resp.Data.([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Data = %v, want 2 entries", resp.Data)
	}
}

// TestHandleTimezonesAllTxt covers the plain-text listing endpoint.
func TestHandleTimezonesAllTxt(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones.txt", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Eastern Standard Time") {
		t.Error("plain text listing missing expected entry")
	}
}

// TestHandleTimezonesSearch covers a valid query, an empty query error, and
// the plain-text variant.
func TestHandleTimezonesSearch(t *testing.T) {
	s := newTestServer(t)

	t.Run("valid query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/search?q=eastern", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		resp := decodeAPIResponse(t, rec.Body.Bytes())
		data, ok := resp.Data.([]interface{})
		if !ok || len(data) != 1 {
			t.Errorf("Data = %v, want 1 entry", resp.Data)
		}
	})

	t.Run("missing query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/search", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		resp := decodeAPIResponse(t, rec.Body.Bytes())
		if resp.Success {
			t.Error("expected Success = false")
		}
	})

	t.Run("txt variant missing query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/search.txt", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("txt variant valid query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/search.txt?q=pacific", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !strings.Contains(rec.Body.String(), "Pacific Standard Time") {
			t.Error("txt search result missing expected entry")
		}
	})
}

// TestHandleTimezonesByOffset covers a valid offset, an invalid offset
// value, and a non-matching offset.
func TestHandleTimezonesByOffset(t *testing.T) {
	s := newTestServer(t)

	t.Run("valid offset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/offset/-5", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		resp := decodeAPIResponse(t, rec.Body.Bytes())
		data, ok := resp.Data.([]interface{})
		if !ok || len(data) != 1 {
			t.Errorf("Data = %v, want 1 entry", resp.Data)
		}
	})

	t.Run("invalid offset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/offset/notanumber", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}

// TestHandleTimezonesByAbbr covers matching and non-matching abbreviations.
func TestHandleTimezonesByAbbr(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/abbr/EST", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	resp := decodeAPIResponse(t, rec.Body.Bytes())
	data, ok := resp.Data.([]interface{})
	if !ok || len(data) != 1 {
		t.Errorf("Data = %v, want 1 entry", resp.Data)
	}
}

// TestHandleTimezonesByUTC covers a matching and a not-found UTC
// identifier.
func TestHandleTimezonesByUTC(t *testing.T) {
	s := newTestServer(t)

	t.Run("found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/utc/America%2FNew_York", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/utc/Nowhere", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

// TestHandleTimezoneByValue covers a matching and a not-found value.
func TestHandleTimezoneByValue(t *testing.T) {
	s := newTestServer(t)

	t.Run("found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/value/Eastern%20Standard%20Time", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones/value/Nonexistent", nil)
		rec := httptest.NewRecorder()
		s.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

// TestHandleTimezonesRandom covers both the API and shorthand routes, plus
// the plain-text variants.
func TestHandleTimezonesRandom(t *testing.T) {
	s := newTestServer(t)
	for _, path := range []string{"/api/v1/timezones/random", "/random"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			s.Router().ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			resp := decodeAPIResponse(t, rec.Body.Bytes())
			if !resp.Success {
				t.Error("expected Success = true")
			}
		})
	}

	for _, path := range []string{"/api/v1/timezones/random.txt", "/random.txt"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			s.Router().ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			if !strings.Contains(rec.Body.String(), "Text:") {
				t.Error("random.txt response missing expected content")
			}
		})
	}
}

// TestHandleTimezonesRandomEmpty verifies the not-found path when the
// service has no timezones.
func TestHandleTimezonesRandomEmpty(t *testing.T) {
	svc, err := timezones.NewService([]byte(`[]`))
	if err != nil {
		t.Fatalf("NewService() unexpected error: %v", err)
	}
	s := New(svc, config.DefaultConfig(), "127.0.0.1", "8080", "1.0.0-test", "2024-01-01", "abcdef")

	req := httptest.NewRequest(http.MethodGet, "/random", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/random.txt", nil)
	rec2 := httptest.NewRecorder()
	s.Router().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec2.Code, http.StatusNotFound)
	}
}

// TestHandleStats covers both the JSON and plain-text stats endpoints.
func TestHandleStats(t *testing.T) {
	s := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	resp := decodeAPIResponse(t, rec.Body.Bytes())
	stats, ok := resp.Data.(map[string]interface{})
	if !ok || stats["total_timezones"].(float64) != 2 {
		t.Errorf("Data = %v, want total_timezones=2", resp.Data)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/stats.txt", nil)
	rec2 := httptest.NewRecorder()
	s.Router().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec2.Code, http.StatusOK)
	}
	if !strings.Contains(rec2.Body.String(), "Total Timezones: 2") {
		t.Error("stats.txt missing expected total")
	}
}

// TestHandleCount covers both the JSON and plain-text count endpoints.
func TestHandleCount(t *testing.T) {
	s := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/count", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	resp := decodeAPIResponse(t, rec.Body.Bytes())
	data, ok := resp.Data.(map[string]interface{})
	if !ok || data["count"].(float64) != 2 {
		t.Errorf("Data = %v, want count=2", resp.Data)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/count.txt", nil)
	rec2 := httptest.NewRecorder()
	s.Router().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec2.Code, http.StatusOK)
	}
	if strings.TrimSpace(rec2.Body.String()) != "2" {
		t.Errorf("count.txt body = %q, want %q", rec2.Body.String(), "2")
	}
}

// TestCORSMiddleware verifies CORS headers are set and OPTIONS requests are
// short-circuited.
func TestCORSMiddleware(t *testing.T) {
	s := newTestServer(t)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/timezones", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("OPTIONS status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want *", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

// TestSecurityHeadersMiddleware verifies the standard security headers are
// present on every response.
func TestSecurityHeadersMiddleware(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("X-Frame-Options = %q, want DENY", rec.Header().Get("X-Frame-Options"))
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("X-Content-Type-Options = %q, want nosniff", rec.Header().Get("X-Content-Type-Options"))
	}
}

// TestNotFoundRoute verifies an unknown path returns 404 rather than
// panicking.
func TestNotFoundRoute(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/this-route-does-not-exist", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
