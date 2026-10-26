package subdomain

import (
	"context"
	"sort"
	"strings"
	"sync"

	"golang.org/x/net/publicsuffix"
)

// Recorder collects resolved hostnames from all passive sources
type Recorder struct {
	mu    sync.Mutex
	hosts map[string]map[string]bool // host -> set of IPs
}

// NewRecorder creates an empty recorder
func NewRecorder() *Recorder {
	return &Recorder{hosts: make(map[string]map[string]bool)}
}

// Add records a hostname and the IPs it resolved to
func (r *Recorder) Add(host string, ips []string) {
	host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
	if host == "" || len(ips) == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	set, ok := r.hosts[host]
	if !ok {
		set = make(map[string]bool)
		r.hosts[host] = set
	}
	for _, ip := range ips {
		set[ip] = true
	}
}

// Results returns host -> sorted IPs
func (r *Recorder) Results() map[string][]string {
	r.mu.Lock()
	defer r.mu.Unlock()

	results := make(map[string][]string, len(r.hosts))
	for host, set := range r.hosts {
		ips := make([]string, 0, len(set))
		for ip := range set {
			ips = append(ips, ip)
		}
		sort.Strings(ips)
		results[host] = ips
	}
	return results
}

type recorderKey struct{}

// WithRecorder returns a context that carries the recorder to passive sources
func WithRecorder(ctx context.Context, r *Recorder) context.Context {
	return context.WithValue(ctx, recorderKey{}, r)
}

// Record adds a resolved host to the context's recorder, if there is one
func Record(ctx context.Context, host string, ips []string) {
	if r, ok := ctx.Value(recorderKey{}).(*Recorder); ok && r != nil {
		r.Add(host, ips)
	}
}

// ApexDomain returns the registrable domain (e.g. api.example.co.uk -> example.co.uk).
// It returns the input unchanged if it cannot be determined.
func ApexDomain(domain string) string {
	domain = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(domain)), ".")
	apex, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return domain
	}
	return apex
}
