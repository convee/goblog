package security

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSessionTokenRoundTrip(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	token, _, err := NewSessionToken("admin@example.com", "password-hash", now)
	if err != nil {
		t.Fatal(err)
	}

	email, err := SessionEmail(token)
	if err != nil || email != "admin@example.com" {
		t.Fatalf("SessionEmail() = %q, %v", email, err)
	}
	email, err = ValidateSessionToken(token, "password-hash", now.Add(time.Hour))
	if err != nil || email != "admin@example.com" {
		t.Fatalf("ValidateSessionToken() = %q, %v", email, err)
	}
}

func TestConfigureSessionSecretRejectsWeakSecret(t *testing.T) {
	if err := ConfigureSessionSecret("too-short"); err == nil {
		t.Fatal("weak session secret was accepted")
	}
}

func TestSessionTokenRejectsTamperingAndExpiry(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	token, _, err := NewSessionToken("admin@example.com", "password-hash", now)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ValidateSessionToken(token+"x", "password-hash", now); err == nil {
		t.Fatal("tampered token was accepted")
	}
	if _, err := ValidateSessionToken(token, "other-hash", now); err == nil {
		t.Fatal("token signed with a different key was accepted")
	}
	if _, err := ValidateSessionToken(token, "password-hash", now.Add(SessionTTL)); err == nil {
		t.Fatal("expired token was accepted")
	}
}

func TestCSRFDoubleSubmit(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://example.com/login", nil)
	token, err := EnsureCSRFToken(recorder, request)
	if err != nil {
		t.Fatal(err)
	}
	response := recorder.Result()
	if len(response.Cookies()) != 1 {
		t.Fatalf("expected one CSRF cookie, got %d", len(response.Cookies()))
	}

	post := httptest.NewRequest("POST", "https://example.com/signin", strings.NewReader(""))
	post.AddCookie(response.Cookies()[0])
	if !ValidateCSRFToken(post, token) {
		t.Fatal("valid CSRF token was rejected")
	}
	if ValidateCSRFToken(post, token+"x") {
		t.Fatal("invalid CSRF token was accepted")
	}
}
