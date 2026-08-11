package security

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
)

var trustedProxies struct {
	sync.RWMutex
	prefixes []netip.Prefix
}

func ConfigureTrustedProxies(cidrs []string) error {
	prefixes := make([]netip.Prefix, 0, len(cidrs))
	for _, cidr := range cidrs {
		prefix, err := netip.ParsePrefix(strings.TrimSpace(cidr))
		if err != nil {
			return fmt.Errorf("invalid trusted proxy %q: %w", cidr, err)
		}
		prefixes = append(prefixes, prefix.Masked())
	}
	trustedProxies.Lock()
	trustedProxies.prefixes = prefixes
	trustedProxies.Unlock()
	return nil
}

// RequestIsHTTPS only trusts X-Forwarded-Proto when the direct peer is a
// configured reverse proxy. Direct TLS requests are always considered secure.
func RequestIsHTTPS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	peer, ok := remoteAddress(r.RemoteAddr)
	if !ok || !isTrustedProxy(peer) {
		return false
	}
	proto := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"), ",")[0])
	return strings.EqualFold(proto, "https")
}

func remoteAddress(remoteAddr string) (netip.Addr, bool) {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	addr, err := netip.ParseAddr(strings.Trim(host, "[]"))
	return addr, err == nil
}

func isTrustedProxy(addr netip.Addr) bool {
	trustedProxies.RLock()
	defer trustedProxies.RUnlock()
	for _, prefix := range trustedProxies.prefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}
