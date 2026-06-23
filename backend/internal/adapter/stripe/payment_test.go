package stripe

import (
	"bytes"
	"context"
	"testing"
	"time"

	stripego "github.com/stripe/stripe-go/v82"
)

// capturingBackend implements stripe.Backend and records the context carried on
// the params of the most recent SDK call, without performing any network I/O.
type capturingBackend struct {
	lastCtx context.Context
}

func (b *capturingBackend) Call(_, _, _ string, params stripego.ParamsContainer, _ stripego.LastResponseSetter) error {
	b.lastCtx = params.GetParams().Context
	return nil
}

func (b *capturingBackend) CallStreaming(_, _, _ string, params stripego.ParamsContainer, _ stripego.StreamingLastResponseSetter) error {
	b.lastCtx = params.GetParams().Context
	return nil
}

func (b *capturingBackend) CallRaw(_, _, _ string, _ []byte, params *stripego.Params, _ stripego.LastResponseSetter) error {
	b.lastCtx = params.Context
	return nil
}

func (b *capturingBackend) CallMultipart(_, _, _, _ string, _ *bytes.Buffer, _ *stripego.Params, _ stripego.LastResponseSetter) error {
	return nil
}

func (b *capturingBackend) SetMaxNetworkRetries(int64) {}

// installCapturingBackend swaps in the fake API backend for the duration of a
// test and restores the original afterwards.
func installCapturingBackend(t *testing.T) *capturingBackend {
	t.Helper()
	fake := &capturingBackend{}
	original := stripego.GetBackend(stripego.APIBackend)
	stripego.SetBackend(stripego.APIBackend, fake)
	t.Cleanup(func() { stripego.SetBackend(stripego.APIBackend, original) })
	return fake
}

func TestPaymentService_CreateCustomer_PassesBoundedContextToSDK(t *testing.T) {
	fake := installCapturingBackend(t)
	svc := NewPaymentServiceWithTimeout("sk_test_x", "whsec_x", 15*time.Second)

	type ctxKey string
	const k ctxKey = "marker"
	parent := context.WithValue(context.Background(), k, "present")

	if _, err := svc.CreateCustomer(parent, "a@b.com", "Alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fake.lastCtx == nil {
		t.Fatal("expected the adapter to set params.Context, got nil")
	}
	// The derived context must still descend from the caller's context.
	if fake.lastCtx.Value(k) != "present" {
		t.Fatal("the SDK context is not derived from the caller's context")
	}
	// And it must carry the per-call deadline ceiling.
	dl, ok := fake.lastCtx.Deadline()
	if !ok {
		t.Fatal("expected a deadline on the SDK context")
	}
	if remaining := time.Until(dl); remaining <= 0 || remaining > 15*time.Second {
		t.Fatalf("deadline out of expected range: %s", remaining)
	}
}

func TestPaymentService_CreateCheckoutSession_PassesContextToSDK(t *testing.T) {
	fake := installCapturingBackend(t)
	svc := NewPaymentServiceWithTimeout("sk_test_x", "whsec_x", 15*time.Second)

	if _, err := svc.CreateCheckoutSession(context.Background(), "cus_1", "price_1", "https://ok", "https://cancel"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.lastCtx == nil {
		t.Fatal("expected params.Context to be set on the checkout call")
	}
	if _, ok := fake.lastCtx.Deadline(); !ok {
		t.Fatal("expected a deadline on the checkout SDK context")
	}
}

func TestPaymentService_RespectsNearerCallerDeadline(t *testing.T) {
	fake := installCapturingBackend(t)
	svc := NewPaymentServiceWithTimeout("sk_test_x", "whsec_x", time.Hour)

	parent, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if _, err := svc.CreateCustomer(parent, "a@b.com", "Alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dl, ok := fake.lastCtx.Deadline()
	if !ok {
		t.Fatal("expected a deadline")
	}
	if remaining := time.Until(dl); remaining > 50*time.Millisecond {
		t.Fatalf("nearer caller deadline should win over the large default, remaining=%s", remaining)
	}
}
