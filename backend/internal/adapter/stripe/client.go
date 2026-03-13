package stripe

import (
	stripego "github.com/stripe/stripe-go/v82"
)

// Init sets the Stripe API key globally.
func Init(apiKey string) {
	stripego.Key = apiKey
}
