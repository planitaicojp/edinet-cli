package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/planitaicojp/edinet-cli/internal/config"
)

var debugEnabled bool

func init() {
	switch os.Getenv(config.EnvDebug) {
	case "1", "true":
		debugEnabled = true
	}
}

// SetDebug enables debug logging. Only enables (never disables).
func SetDebug(enabled bool) {
	if enabled {
		debugEnabled = true
	}
}

func debugLogRequest(req *http.Request) {
	if !debugEnabled {
		return
	}
	// Mask Subscription-Key in URL
	maskedURL := maskSubscriptionKey(req.URL.String())
	fmt.Fprintf(os.Stderr, "> %s %s\n", req.Method, maskedURL)
	for name, values := range req.Header {
		fmt.Fprintf(os.Stderr, "> %s: %s\n", name, strings.Join(values, ", "))
	}
	fmt.Fprintln(os.Stderr)
}

func debugLogResponse(resp *http.Response, duration time.Duration, body []byte) {
	if !debugEnabled {
		return
	}
	fmt.Fprintf(os.Stderr, "< %d %s (%dms)\n", resp.StatusCode, http.StatusText(resp.StatusCode), duration.Milliseconds())
	for name, values := range resp.Header {
		fmt.Fprintf(os.Stderr, "< %s: %s\n", name, strings.Join(values, ", "))
	}
	if len(body) > 0 {
		s := string(body)
		if len(s) > 2000 {
			s = s[:2000] + "...(truncated)"
		}
		fmt.Fprintf(os.Stderr, "< %s\n", s)
	}
	fmt.Fprintln(os.Stderr)
}

func debugLogRetry(attempt int, statusCode int, wait time.Duration) {
	if !debugEnabled {
		return
	}
	fmt.Fprintf(os.Stderr, "* retry %d/%d (HTTP %d), waiting %s\n", attempt, maxRetries, statusCode, wait)
}

// maskSubscriptionKey replaces the Subscription-Key query parameter value with a masked version.
func maskSubscriptionKey(rawURL string) string {
	const key = "Subscription-Key="
	idx := strings.Index(rawURL, key)
	if idx < 0 {
		return rawURL
	}
	start := idx + len(key)
	end := strings.IndexByte(rawURL[start:], '&')
	if end < 0 {
		return rawURL[:start] + "****"
	}
	return rawURL[:start] + "****" + rawURL[start+end:]
}
