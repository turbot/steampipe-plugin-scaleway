package scaleway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// serveTooManyRequests answers 429 to the first `rejections` requests, then 200.
func serveTooManyRequests(t *testing.T, rejections int32, retryAfter string) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) <= rejections {
			if retryAfter != "" {
				w.Header().Set("Retry-After", retryAfter)
			}
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	return server, &calls
}

func fastTransport() *rateLimitRetryTransport {
	transport := newRateLimitRetryTransport(http.DefaultTransport)
	transport.baseDelay = time.Millisecond
	transport.maxDelay = 5 * time.Millisecond
	return transport
}

func TestRateLimitRetryTransportRetriesUntilSuccess(t *testing.T) {
	server, calls := serveTooManyRequests(t, 3, "")
	client := &http.Client{Transport: fastTransport()}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || calls.Load() != 4 {
		t.Fatalf("got status %d after %d calls, want 200 after 4", resp.StatusCode, calls.Load())
	}
}

func TestRateLimitRetryTransportGivesUpAfterMaxRetries(t *testing.T) {
	server, calls := serveTooManyRequests(t, 100, "")
	transport := fastTransport()
	client := &http.Client{Transport: transport}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// The last 429 reaches the caller, which reports it as before
	want := int32(transport.maxRetries + 1)
	if resp.StatusCode != http.StatusTooManyRequests || calls.Load() != want {
		t.Fatalf("got status %d after %d calls, want 429 after %d", resp.StatusCode, calls.Load(), want)
	}
}

func TestRateLimitRetryTransportDoesNotReplayWrites(t *testing.T) {
	server, calls := serveTooManyRequests(t, 1, "")
	client := &http.Client{Transport: fastTransport()}

	resp, err := client.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests || calls.Load() != 1 {
		t.Fatalf("got status %d after %d calls, want 429 after 1", resp.StatusCode, calls.Load())
	}
}

func TestRateLimitRetryTransportStopsWhenTheContextEnds(t *testing.T) {
	server, _ := serveTooManyRequests(t, 100, "")
	transport := newRateLimitRetryTransport(http.DefaultTransport)
	client := &http.Client{Transport: transport}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)

	start := time.Now()
	if resp, err := client.Do(req); err == nil {
		resp.Body.Close()
		t.Fatal("want an error once the context is done")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("waited %s, want the retry to stop with the context", elapsed)
	}
}

func TestRateLimitRetryTransportDelay(t *testing.T) {
	transport := newRateLimitRetryTransport(http.DefaultTransport)

	if got := transport.delay(0, "3"); got != 3*time.Second {
		t.Errorf("Retry-After 3: got %s, want 3s", got)
	}
	if got := transport.delay(0, "3600"); got != transport.maxDelay {
		t.Errorf("Retry-After 3600: got %s, want the %s cap", got, transport.maxDelay)
	}
	for attempt, upper := range []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 10 * time.Second, 10 * time.Second} {
		if got := transport.delay(attempt, ""); got < upper/2 || got > upper {
			t.Errorf("attempt %d: got %s, want between %s and %s", attempt, got, upper/2, upper)
		}
	}
}
