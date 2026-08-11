package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	SessionCookieName = "goblog_session"
	CSRFCookieName    = "goblog_csrf"
	SessionTTL        = 12 * time.Hour
)

var errInvalidSession = errors.New("invalid session")

var (
	sessionSecretMu sync.RWMutex
	sessionSecret   = newRandomSecret()
)

type sessionClaims struct {
	Email   string `json:"email"`
	Expires int64  `json:"expires"`
}

// NewSessionToken creates an expiring token signed with the application secret
// and the user's password hash. Changing either invalidates existing sessions.
func NewSessionToken(email, passwordHash string, now time.Time) (string, time.Time, error) {
	if email == "" || passwordHash == "" {
		return "", time.Time{}, errInvalidSession
	}

	expires := now.Add(SessionTTL)
	claims, err := json.Marshal(sessionClaims{Email: email, Expires: expires.Unix()})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("encode session: %w", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(claims)
	signature := sign(payload, passwordHash)
	return payload + "." + signature, expires, nil
}

// ConfigureSessionSecret replaces the process-local fallback secret. Production
// deployments must call this with a stable secret from an external secret store.
func ConfigureSessionSecret(secret string) error {
	if len(secret) < 32 {
		return errors.New("session secret must contain at least 32 bytes")
	}
	sessionSecretMu.Lock()
	sessionSecret = []byte(secret)
	sessionSecretMu.Unlock()
	return nil
}

// SessionEmail extracts the untrusted email claim so the caller can load the
// password hash needed to validate the signature.
func SessionEmail(token string) (string, error) {
	claims, _, err := parseSessionToken(token)
	if err != nil || claims.Email == "" {
		return "", errInvalidSession
	}
	return claims.Email, nil
}

// ValidateSessionToken verifies the signature and expiry and returns the
// authenticated email address.
func ValidateSessionToken(token, passwordHash string, now time.Time) (string, error) {
	claims, signature, err := parseSessionToken(token)
	if err != nil || passwordHash == "" || claims.Email == "" {
		return "", errInvalidSession
	}
	payload := strings.SplitN(token, ".", 2)[0]
	expected, err := base64.RawURLEncoding.DecodeString(sign(payload, passwordHash))
	if err != nil || !hmac.Equal(signature, expected) || now.Unix() >= claims.Expires {
		return "", errInvalidSession
	}
	return claims.Email, nil
}

func parseSessionToken(token string) (sessionClaims, []byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return sessionClaims{}, nil, errInvalidSession
	}
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return sessionClaims{}, nil, errInvalidSession
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || len(signature) != sha256.Size {
		return sessionClaims{}, nil, errInvalidSession
	}
	var claims sessionClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return sessionClaims{}, nil, errInvalidSession
	}
	return claims, signature, nil
}

func sign(payload, key string) string {
	sessionSecretMu.RLock()
	secret := append([]byte(nil), sessionSecret...)
	sessionSecretMu.RUnlock()
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(payload))
	_, _ = mac.Write([]byte{0})
	_, _ = mac.Write([]byte(key))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func newRandomSecret() []byte {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		panic(fmt.Sprintf("generate session secret: %v", err))
	}
	return secret
}

func SetSessionCookie(w http.ResponseWriter, r *http.Request, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		MaxAge:   int(SessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   RequestIsHTTPS(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearSessionCookies(w http.ResponseWriter, r *http.Request) {
	for _, name := range []string{SessionCookieName, CSRFCookieName, "email"} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   RequestIsHTTPS(r),
			SameSite: http.SameSiteLaxMode,
		})
	}
}

// EnsureCSRFToken returns the existing token or issues a new host-only token.
func EnsureCSRFToken(w http.ResponseWriter, r *http.Request) (string, error) {
	if cookie, err := r.Cookie(CSRFCookieName); err == nil && validCSRFToken(cookie.Value) {
		return cookie.Value, nil
	}

	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("generate csrf token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(random)
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(SessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   RequestIsHTTPS(r),
		SameSite: http.SameSiteStrictMode,
	})
	return token, nil
}

func ValidateCSRFToken(r *http.Request, submitted string) bool {
	if !validCSRFToken(submitted) {
		return false
	}
	cookie, err := r.Cookie(CSRFCookieName)
	if err != nil || !validCSRFToken(cookie.Value) || len(cookie.Value) != len(submitted) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(submitted)) == 1
}

func validCSRFToken(token string) bool {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	return err == nil && len(decoded) == 32
}
