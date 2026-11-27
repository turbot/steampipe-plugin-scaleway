package scaleway

import (
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// rateLimitRetryTransport replays a request the Scaleway API rejected with
// 429 Too Many Requests.
//
// The scaleway-sdk-go client does not retry at all, and the Steampipe SDK
// cannot help either: it refuses to retry a list hydrate once a row has been
// streamed, so a 429 on a later page, or in a later zone or region of the
// matrix, would fail the whole table. Retrying each HTTP request covers every
// call, mid-pagination included. A table that fans out over every zone sends
// many requests at once, which is exactly when the API starts answering 429.
type rateLimitRetryTransport struct {
	next       http.RoundTripper
	maxRetries int
	baseDelay  time.Duration
	maxDelay   time.Duration
}

func newRateLimitRetryTransport(next http.RoundTripper) *rateLimitRetryTransport {
	return &rateLimitRetryTransport{
		next:       next,
		maxRetries: 5,
		baseDelay:  time.Second,
		maxDelay:   10 * time.Second,
	}
}

func (t *rateLimitRetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only replay requests that are safe to send twice and carry no body
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		return t.next.RoundTrip(req)
	}

	for attempt := 0; ; attempt++ {
		resp, err := t.next.RoundTrip(req)
		if err != nil || resp.StatusCode != http.StatusTooManyRequests || attempt >= t.maxRetries {
			return resp, err
		}

		delay := t.delay(attempt, resp.Header.Get("Retry-After"))
		// Drain the body so the connection can be reused
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()

		timer := time.NewTimer(delay)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}
}

// delay returns how long to wait before the next attempt: the Retry-After the
// API asked for, in seconds, else an exponential backoff, both capped. The
// backoff is jittered over its upper half so that the many connections hitting
// the same limit at the same time do not all come back together.
func (t *rateLimitRetryTransport) delay(attempt int, retryAfter string) time.Duration {
	if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds >= 0 {
		return min(time.Duration(seconds)*time.Second, t.maxDelay)
	}
	d := min(t.baseDelay<<attempt, t.maxDelay)
	return d/2 + rand.N(d/2+1)
}
