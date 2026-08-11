package security

import (
	"crypto/tls"
	"net/http/httptest"
	"testing"
)

func TestRequestIsHTTPS(t *testing.T) {
	t.Cleanup(func() { _ = ConfigureTrustedProxies(nil) })

	directTLS := httptest.NewRequest("GET", "https://example.com", nil)
	directTLS.TLS = &tls.ConnectionState{}
	if !RequestIsHTTPS(directTLS) {
		t.Fatal("direct TLS request should be secure")
	}

	untrusted := httptest.NewRequest("GET", "http://example.com", nil)
	untrusted.RemoteAddr = "203.0.113.10:1234"
	untrusted.Header.Set("X-Forwarded-Proto", "https")
	if RequestIsHTTPS(untrusted) {
		t.Fatal("untrusted forwarded header must be ignored")
	}

	if err := ConfigureTrustedProxies([]string{"10.0.0.0/8"}); err != nil {
		t.Fatal(err)
	}
	trusted := httptest.NewRequest("GET", "http://example.com", nil)
	trusted.RemoteAddr = "10.1.2.3:4321"
	trusted.Header.Set("X-Forwarded-Proto", "https")
	if !RequestIsHTTPS(trusted) {
		t.Fatal("trusted proxy header should be honored")
	}
}

func TestConfigureTrustedProxiesRejectsInvalidCIDR(t *testing.T) {
	if err := ConfigureTrustedProxies([]string{"not-a-cidr"}); err == nil {
		t.Fatal("expected invalid CIDR error")
	}
}
