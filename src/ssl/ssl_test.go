package ssl

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// generateSelfSignedPEM produces a minimal self-signed certificate/key pair
// (PEM-encoded) for the given domain, suitable for exercising
// tls.LoadX509KeyPair without any real CA or network access.
func generateSelfSignedPEM(t *testing.T, domain string) (certPEM, keyPEM []byte) {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: domain},
		DNSNames:     []string{domain},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:         true,
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("CreateCertificate: %v", err)
	}

	var certBuf bytes.Buffer
	if err := pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		t.Fatalf("pem.Encode cert: %v", err)
	}

	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatalf("MarshalECPrivateKey: %v", err)
	}
	var keyBuf bytes.Buffer
	if err := pem.Encode(&keyBuf, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}); err != nil {
		t.Fatalf("pem.Encode key: %v", err)
	}

	return certBuf.Bytes(), keyBuf.Bytes()
}

// TestParseChallenge covers every recognized alias, case-insensitivity, and
// the fallback-to-http-01 default.
func TestParseChallenge(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "http-01", input: "http-01", want: "http-01"},
		{name: "http01", input: "http01", want: "http-01"},
		{name: "http", input: "http", want: "http-01"},
		{name: "tls-alpn-01", input: "tls-alpn-01", want: "tls-alpn-01"},
		{name: "tlsalpn01", input: "tlsalpn01", want: "tls-alpn-01"},
		{name: "tls-alpn", input: "tls-alpn", want: "tls-alpn-01"},
		{name: "tls", input: "tls", want: "tls-alpn-01"},
		{name: "dns-01", input: "dns-01", want: "dns-01"},
		{name: "dns01", input: "dns01", want: "dns-01"},
		{name: "dns", input: "dns", want: "dns-01"},
		{name: "uppercase and whitespace", input: "  HTTP-01  ", want: "http-01"},
		{name: "unknown falls back to http-01", input: "carrier-pigeon", want: "http-01"},
		{name: "empty falls back to http-01", input: "", want: "http-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseChallenge(tt.input); got != tt.want {
				t.Errorf("ParseChallenge(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestChallengeServer covers setting, clearing, and serving tokens, plus the
// non-ACME-path passthrough (returns false, does not write a response).
func TestChallengeServer(t *testing.T) {
	cs := NewChallengeServer()

	t.Run("non-acme path is not handled", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/some/other/path", nil)
		rec := httptest.NewRecorder()
		if handled := cs.ServeHTTP(rec, req); handled {
			t.Error("ServeHTTP() = true for non-acme path, want false")
		}
	})

	t.Run("unknown token returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/unknown-token", nil)
		rec := httptest.NewRecorder()
		if handled := cs.ServeHTTP(rec, req); !handled {
			t.Fatal("ServeHTTP() = false for acme path, want true")
		}
		if rec.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("known token is served", func(t *testing.T) {
		cs.SetToken("tok1", "tok1.keyauth")
		req := httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/tok1", nil)
		rec := httptest.NewRecorder()
		if handled := cs.ServeHTTP(rec, req); !handled {
			t.Fatal("ServeHTTP() = false for known token, want true")
		}
		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if rec.Body.String() != "tok1.keyauth" {
			t.Errorf("body = %q, want %q", rec.Body.String(), "tok1.keyauth")
		}

		cs.ClearToken("tok1")
		req2 := httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/tok1", nil)
		rec2 := httptest.NewRecorder()
		cs.ServeHTTP(rec2, req2)
		if rec2.Code != http.StatusNotFound {
			t.Errorf("status after ClearToken = %d, want %d", rec2.Code, http.StatusNotFound)
		}
	})
}

// TestNewManagerDisabled verifies GetTLSConfig short-circuits to nil, nil
// when SSL is not enabled, regardless of certificate availability.
func TestNewManagerDisabled(t *testing.T) {
	m := NewManager(Config{Enabled: false})
	cfg, err := m.GetTLSConfig([]string{"example.com"})
	if err != nil {
		t.Fatalf("GetTLSConfig() unexpected error: %v", err)
	}
	if cfg != nil {
		t.Errorf("GetTLSConfig() = %v, want nil when disabled", cfg)
	}
}

// TestGetTLSConfigNoCerts verifies an informative error is returned when SSL
// is enabled but no certificates exist and Let's Encrypt is off.
func TestGetTLSConfigNoCerts(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(Config{Enabled: true, CertPath: dir})
	_, err := m.GetTLSConfig([]string{"example.com"})
	if err == nil {
		t.Fatal("GetTLSConfig() expected error when no certs and Let's Encrypt disabled")
	}
}

// TestGetTLSConfigManualCerts verifies manually placed cert/key files (both
// naming conventions) are found and loaded.
func TestGetTLSConfigManualCerts(t *testing.T) {
	dir := t.TempDir()
	domain := "example.com"

	certPEM, keyPEM := generateSelfSignedPEM(t, domain)

	t.Run("flat .crt/.key naming", func(t *testing.T) {
		sub := filepath.Join(dir, "flat")
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(filepath.Join(sub, domain+".crt"), certPEM, 0644); err != nil {
			t.Fatalf("WriteFile cert: %v", err)
		}
		if err := os.WriteFile(filepath.Join(sub, domain+".key"), keyPEM, 0600); err != nil {
			t.Fatalf("WriteFile key: %v", err)
		}

		m := NewManager(Config{Enabled: true, CertPath: sub})
		tlsCfg, err := m.GetTLSConfig([]string{domain})
		if err != nil {
			t.Fatalf("GetTLSConfig() unexpected error: %v", err)
		}
		if tlsCfg == nil || len(tlsCfg.Certificates) != 1 {
			t.Errorf("GetTLSConfig() = %v, want a config with 1 certificate", tlsCfg)
		}
	})

	t.Run("fullchain/privkey directory naming", func(t *testing.T) {
		sub := filepath.Join(dir, "fullchain")
		domDir := filepath.Join(sub, domain)
		if err := os.MkdirAll(domDir, 0755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(filepath.Join(domDir, "fullchain.pem"), certPEM, 0644); err != nil {
			t.Fatalf("WriteFile cert: %v", err)
		}
		if err := os.WriteFile(filepath.Join(domDir, "privkey.pem"), keyPEM, 0600); err != nil {
			t.Fatalf("WriteFile key: %v", err)
		}

		m := NewManager(Config{Enabled: true, CertPath: sub})
		tlsCfg, err := m.GetTLSConfig([]string{domain})
		if err != nil {
			t.Fatalf("GetTLSConfig() unexpected error: %v", err)
		}
		if tlsCfg == nil || len(tlsCfg.Certificates) != 1 {
			t.Errorf("GetTLSConfig() = %v, want a config with 1 certificate", tlsCfg)
		}
	})
}

// TestGetHTTPHandlerFallback verifies that without a certManager configured
// (i.e. Let's Encrypt was never engaged), GetHTTPHandler simply returns the
// fallback handler unchanged.
func TestGetHTTPHandlerFallback(t *testing.T) {
	m := NewManager(Config{Enabled: false})
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler := m.GetHTTPHandler(fallback)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want %d (fallback should be used untouched)", rec.Code, http.StatusTeapot)
	}
}

// TestFileExists covers an existing file and a nonexistent path.
func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "present.txt")
	if err := os.WriteFile(existing, []byte("x"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if !fileExists(existing) {
		t.Error("fileExists() = false for existing file, want true")
	}
	if fileExists(filepath.Join(dir, "missing.txt")) {
		t.Error("fileExists() = true for missing file, want false")
	}
}
